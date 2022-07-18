package redis_helper

import (
	"context"
	"entry-task/app_data_model"
	"entry-task/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var ctx = context.Background()

var rdb *redis.Client

func InitRedisClient(config config.RedisConfigs) {
	addr := config.Host + ":" + strconv.Itoa(config.Port)

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})
}

func GetInfoFromRedisById(sessionId string) (app_data_model.UserInfo, bool) {
	fmt.Println("Try to get info from redis, session id: ", sessionId)
	res, _ := rdb.LRange(ctx, sessionId, 0, -1).Result()

	if len(res) == 0 {
		fmt.Println("No such key in redis")
		return app_data_model.UserInfo{}, false
	}

	fmt.Println("Get info from redis successfully, res: ", res)
	return app_data_model.UserInfo{res[2], res[1], res[0]}, true
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SetDataToRedisById(sessionId string, userInfo app_data_model.UserInfo) {
	fmt.Println("Set data into redis, user info: ", userInfo)
	keys := rdb.Keys(ctx, "*").Val()
	if contains(keys, sessionId) {
		rdb.Del(ctx, sessionId)
	}
	rdb.LPush(ctx, sessionId, userInfo.Nickname, userInfo.Password, userInfo.PhotoPath)
}

func UpdatePhotoPathById(sessionId string, photoPath string) {
	rdb.LSet(ctx, sessionId, 0, photoPath)
}

func UpdateNicknameById(sessionId string, nickName string) {
	rdb.LSet(ctx, sessionId, 2, nickName)
}
