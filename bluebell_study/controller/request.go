package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUserID GetCurrentUser 获取当前登录的用户
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey) //这里是通过grom获取用户id
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (page, size int64, err error) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			page = 1
		}
	}

	if sizeStr == "" {
		size = 10
	} else {
		size, err = strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			size = 10
		}
	}

	return page, size, nil
}
