package controller

import (
	"bluebell_study/logic"
	"bluebell_study/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//投票

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)

	// 先打印原始请求体
	zap.L().Info("收到投票请求",
		zap.String("ContentType", c.ContentType()),
		zap.Int64("ContentLength", c.Request.ContentLength))

	if err := c.ShouldBindJSON(p); err != nil {
		// 打印详细的错误信息
		zap.L().Error("ShouldBindJSON failed",
			zap.Error(err),
			zap.String("error_type", fmt.Sprintf("%T", err)),
			zap.Any("param_vote_data", p))

		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			zap.L().Error("类型断言失败，不是ValidationErrors",
				zap.String("actual_error_type", fmt.Sprintf("%T", err)))
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	zap.L().Info("参数绑定成功", zap.Any("param_vote_data", p))

	// 获取当前请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
