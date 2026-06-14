package middlewares

import (
	"bluebell_study/controller"
	"bluebell_study/dao/redis"
	"bluebell_study/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		//这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 请求头格式：Bearer xxxxxxx.xxx.xxx
		//这里的具体实现方式主要依据实际业务需求
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		//按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			zap.L().Error("jwt.ParseToken failed",
				zap.String("token", parts[1]),
				zap.Error(err))
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 验证Redis中的Token是否匹配（单设备登录检查）
		zap.L().Debug("checking redis token", zap.Int64("userID", mc.UserID))
		storedToken, err := redis.GetUserToken(mc.UserID)
		if err != nil {
			// Redis中没有找到Token，说明用户未登录或Token已失效
			zap.L().Error("redis.GetUserToken failed",
				zap.Int64("userID", mc.UserID),
				zap.Error(err))
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		if storedToken != parts[1] {
			// Token不匹配，说明在其他设备登录了
			zap.L().Warn("token mismatch",
				zap.Int64("userID", mc.UserID),
				zap.String("storedToken", storedToken[:20]+"..."),
				zap.String("requestToken", parts[1][:20]+"..."))
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		zap.L().Info("auth success", zap.Int64("userID", mc.UserID))
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理请求函数中，可以用过c.Get(CtxUserIDKey) 来获取当前请求的用户信息
	}
}
