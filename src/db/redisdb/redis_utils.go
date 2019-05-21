package redisdb

import (
	"strconv"
	"strings"
)

func (r *Redis) plusyOffsetKey() string {
	return "plusy:offset"
}

func (r *Redis) chatIdAsString(chatId int64) string {
	return strconv.FormatInt(chatId, 10)
}

func (r *Redis) chatKeysKey(chatId int64) string {
	return "chat:keys:" + strconv.FormatInt(chatId, 10)
}

func (r *Redis) chatKeyPrefix(chatId int64) string {
	return "chat:" + strconv.FormatInt(chatId, 10)
}

func (r *Redis) userIdAsInt64(userIdAsString string) int64 {
	userId, err := strconv.ParseInt(userIdAsString, 10, 64)
	if err != nil {
		return 0
	}

	return userId
}

func (r *Redis) userIdToString(userId int64) string {
	return strconv.FormatInt(userId, 10)
}

func (r *Redis) userIdKey(userId int64) string {
	return "user:id:" + strconv.FormatInt(userId, 10)
}

func (r *Redis) userUsernameKey(username string) string {
	return "user:un:" + username
}

func (r *Redis) userFirstnameKey(firstname string) string {
	// lowercase firstname before storing so users can query case insensitively
	return "user:fn:" + strings.ToLower(firstname)
}
