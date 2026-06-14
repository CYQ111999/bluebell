package logic

import (
	"bluebell_study/dao/redis"
	"bluebell_study/models"
	"strconv"

	"go.uber.org/zap"
)

// 投票功能
// 1.用户投票

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Any("userID", userID), zap.Any("paramVoteData", p))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
