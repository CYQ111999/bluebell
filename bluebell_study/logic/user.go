package logic

import (
	"bluebell_study/dao/mysql"
	"bluebell_study/dao/redis"
	"bluebell_study/models"
	"bluebell_study/pkg/jwt"
	"bluebell_study/pkg/snowflake"
	"bluebell_study/setting"
	"time"
)

// SignUp 注册逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户存不存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错

		return err
	}
	//生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	u := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//保存进数据库
	return mysql.InsertUser(&u)
}

// Login 登录逻辑
func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 数据库查询,传递指针
	if err := mysql.Login(user); err != nil {
		return "", err
	}
	// 生成JWT
	token, err = jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return "", err
	}

	// 将Token存入Redis，实现单设备登录
	expireHours := setting.Conf.JWTExpire
	err = redis.SetUserToken(user.UserID, token, time.Duration(expireHours)*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}
