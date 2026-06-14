package logic

import (
	"bluebell_study/dao/mysql"
	"bluebell_study/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查到所有community并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
