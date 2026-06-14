package controller

import (
	"bluebell_study/dao/mysql"
	"bluebell_study/dao/redis"
	"bluebell_study/logic"
	"bluebell_study/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求函数
func SignUpHandler(c *gin.Context) {
	//1. 参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误直接返回响应
		zap.L().Error("SignU with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	//1. 参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误直接返回响应
		zap.L().Error("SignU with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2. 业务处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("login.Login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	//3. 返回响应
	ResponseSuccess(c, token)
}

// LogoutHandler 处理登出请求
// 这里的函数暂时是为了实现同账号单设备登录，利用的是redis里存储并比对同一id的token是否不同，不同就登出。
// 因此需要实现登出前删除自己的token的操作
func LogoutHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get(CtxUserIDKey)
	if !exists {
		ResponseError(c, CodeNeedLogin)
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		ResponseError(c, CodeOtherOnline)
		return
	}

	// 从Redis中删除Token
	if err := redis.DeleteUserToken(uid); err != nil {
		zap.L().Error("logout failed", zap.Error(err))
		ResponseError(c, CodeOtherOnline)
		return
	}

	ResponseSuccess(c, nil)
}
