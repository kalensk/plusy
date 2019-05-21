package redisdb

import (
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/kalensk/plusy/src/messages"
	"github.com/kalensk/plusy/src/utils"
	"github.com/pkg/errors"
)

func (r *Redis) IncrementCount(chatId int64, dateGiven int64, giverUser *messages.User, receiverUser *messages.User) (string, error) {
	r.log.Debugf("Incrementing count for user: %v", receiverUser)

	isBot, err := r.isBot(giverUser)
	if err != nil {
		return "", errors.Wrapf(err, "failed checking if user '%+v' is a bot", giverUser)
	}

	if isBot {
		return "", errors.Errorf("failed to give point to bot %v in chat %s (%d)", giverUser, "friendlyName", chatId)
	}

	isBot, err = r.isBot(receiverUser)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check if user '%+v' is valid", receiverUser)
	}

	if isBot {
		return "", errors.Errorf("failed to give point to invalid user %v in chat %s (%d)", receiverUser, "friendlyName", chatId)
	}

	spo := fmt.Sprintf("spo:%d:gave:%d", giverUser.ID, receiverUser.ID)
	spoKey := "chat:" + r.chatIdAsString(chatId) + ":" + spo

	ops := fmt.Sprintf("ops:%d:gave:%d", receiverUser.ID, giverUser.ID)
	opsKey := "chat:" + r.chatIdAsString(chatId) + ":" + ops

	// TODO: figure out how to handle errors inside the transaction
	// https://redis.io/topics/transactions

	// GET do a multi and exec since getUser in log saveTimestamp does a db call for the user
	_, err = r.conn.Do("MULTI")
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'MULTI' failed")
	}

	// chat:keys lexicographical set
	_, err = r.conn.Do("ZADD", "chat:keys:"+r.chatIdAsString(chatId), 0, spo)
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'ZADD chat:keys:%s 0 %s' failed", r.chatIdAsString(chatId), spo)
	}

	_, err = r.conn.Do("ZADD", "chat:keys:"+r.chatIdAsString(chatId), 0, ops)
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'ZADD chat:keys:%s 0 %s' failed", r.chatIdAsString(chatId), ops)
	}

	// spo timestamp
	err = r.saveTimeStamp(chatId, dateGiven, giverUser, receiverUser)
	if err != nil {
		// ToDo: create a helper method to get friendly chat room from chatId. For log messages.
		return "", errors.Wrapf(err, "failed to save timestamp for chat %s (%d) on %s to %+v from %+v",
			"friendlyName", chatId, dateGiven, giverUser, receiverUser)
	}

	// chat:top sorted set
	_, err = r.conn.Do("ZADD", "chat:top:"+r.chatIdAsString(chatId), "INCR", 1, receiverUser.ID)
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'ZADD chat:top:%s INCR 1 %s' failed", r.chatIdAsString(chatId), receiverUser.ID)
	}
	_, err = r.conn.Do("ZREMRANGEBYRANK", "chat:top:"+r.chatIdAsString(chatId), 0, -11) // keep only the top 10
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'ZREMRANGEBYRANK chat:top:%s 0 -11' failed", r.chatIdAsString(chatId))
	}

	// spo:score
	_, err = redis.String(r.conn.Do("INCR", spoKey+":score"))
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'INCR %s:score' failed ", spoKey)
	}

	_, err = redis.String(r.conn.Do("INCR", opsKey+":score"))
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'INCR %s:score' failed ", opsKey)
	}

	// EXEC cannot be deferred since we want to return the result of INCR, which is the currentCount.
	batchedReplies, err := redis.Values(r.conn.Do("EXEC"))
	if err != nil {
		return "", errors.Wrapf(err, "redis command 'EXEC' failed")
	}

	currentCount := batchedReplies[len(batchedReplies)-2].(int64)
	return strconv.FormatInt(currentCount, 10), nil
}

func (r *Redis) saveTimeStamp(chatId int64, dateGiven int64, giverUser *messages.User, receiverUser *messages.User) error {
	// need to defer the getUser() lookup since it is within a redis transaction (MULTI/EXEC)
	r.log.Debugf("Saving timestamp for chat %s (%d) on %d to user %+v from user %+v", "friendlyChatname", chatId, dateGiven, giverUser, receiverUser)

	spo := fmt.Sprintf("spo:%d:gave:%d", giverUser.ID, receiverUser.ID)
	spoKey := "chat:" + r.chatIdAsString(chatId) + ":" + spo

	ops := fmt.Sprintf("ops:%d:gave:%d", receiverUser.ID, giverUser.ID)
	opsKey := "chat:" + r.chatIdAsString(chatId) + ":" + ops

	_, err := r.conn.Do("LPUSH", spoKey+":ts", dateGiven)
	if err != nil {
		return errors.Wrapf(err, "redis command 'LPUSH %s:ts %s' failed", spoKey, dateGiven)
	}
	_, err = r.conn.Do("LPUSH", opsKey+":ts", dateGiven)
	if err != nil {
		return errors.Wrapf(err, "redis command 'LPUSH %s:ts %s' failed", opsKey, dateGiven)
	}

	_, err = r.conn.Do("LTRIM", spoKey+":ts", 0, 500) // ToDo: extract 500 (max size of point timestamps) into a variable
	if err != nil {
		return errors.Wrapf(err, "redis command 'LTRIM %s:ts 0 500' failed", spoKey)
	}

	_, err = r.conn.Do("LTRIM", opsKey+":ts", 0, 500) // ToDo: extract 500 (max size of point timestamps) into a variable
	if err != nil {
		return errors.Wrapf(err, "redis command 'LTRIM %s:ts 0 500' failed", opsKey)
	}

	return nil
}

// /stats tokie
// tokie has 2 plusies

// Received By:
//   doug 1
//   brian 1
//
// Given 1 plusy:
//   kaboodle 1

func (r *Redis) GetPointsReceived(chatId int64, giverUserId int64) ([]messages.UserPoints, error) {
	// ToDo: Multi and Exec
	giverUser, err := r.getUser(giverUserId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user with userId '%d'", giverUserId)
	}

	r.log.Debugf("Getting points received for %v in chat %v", giverUser, chatId)
	// ops:z:gave:a  == z:received:a == opposite of spo:a:gave:z
	opsKey := "ops:" + r.userIdToString(giverUserId) + ":gave:"
	return r.getUserPoints(opsKey, chatId, giverUserId)
}

func (r *Redis) GetPointsGiven(chatId int64, receiverUserId int64) ([]messages.UserPoints, error) {
	// ToDo: Multi and Exec
	receiverUser, err := r.getUser(receiverUserId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user with userId '%d'", receiverUserId)
	}

	r.log.Debugf("Getting points given for %v in chat &v", receiverUser, chatId)
	spoKey := "spo:" + r.userIdToString(receiverUserId) + ":gave:"
	return r.getUserPoints(spoKey, chatId, receiverUserId)
}

func (r *Redis) getUserPoints(lexicographicalKey string, chatId int64, userId int64) ([]messages.UserPoints, error) {
	// ToDo: Multi and Exec
	user, err := r.getUser(userId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user with userId '%d'", userId)
	}

	r.log.Debugf("Getting user points for %v in chat &v", user, chatId)

	users, err := redis.Strings(r.conn.Do("ZRANGEBYLEX", r.chatKeysKey(chatId), "["+lexicographicalKey, "["+lexicographicalKey+"\xFF"))
	if err != nil {
		return nil, errors.Wrapf(err, "redis command 'ZRANGEBYLEX %s [%s [%s\xFF' failed", r.chatKeysKey(chatId), lexicographicalKey, lexicographicalKey)
	}

	var userPoints []messages.UserPoints
	for _, key := range users {
		scoreKey := r.chatKeyPrefix(chatId) + ":" + key + ":score"
		plusies, err := redis.String(r.conn.Do("GET", scoreKey))
		if err != nil {
			return nil, errors.Wrapf(err, "redis command 'GET %s' failed", scoreKey)
		}

		keys := strings.Split(key, ":")
		givenUserId := keys[1]

		givenUser, err := r.getUser(r.userIdAsInt64(givenUserId))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get user with userId '%d'", r.userIdAsInt64(givenUserId))
		}

		userPoints = append(userPoints, messages.UserPoints{User: givenUser, Plusies: plusies})
	}

	return userPoints, nil
}

func (r *Redis) GetTopN(chatId int64, n int) ([]messages.UserPoints, error) {
	key := "chat:top:" + r.chatIdAsString(chatId)
	values, err := redis.StringMap(r.conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES", "LIMIT", 0, n))
	if err != nil {
		return nil, errors.Wrapf(err, "redis command 'ZREVRANGEBYSCORE %s +inf -inf WITHSCORES LIMIT 0 %d' failed", key, n)
	}

	var topN []messages.UserPoints
	for userIdAsString, points := range values {
		user, err := r.getUser(r.userIdAsInt64(userIdAsString))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get user with userId '%d'", r.userIdAsInt64(userIdAsString))
		}

		topN = append(topN, messages.UserPoints{User: user, Plusies: points})
	}

	return topN, nil
}

// ToDo: What if user rename's.... Is there a rename event?
func (r *Redis) SaveOrRemoveUser(message *messages.Message) error {
	user, err := r.getUser(message.From.ID)
	if err != nil && err != redis.ErrNil {
		return errors.Wrapf(err, "failed to get user '%+v'", message.From)
	}

	if user == nil {
		err = r.saveUser(message.From)
	}

	if message.NewChatMember != nil {
		err = r.saveUser(message.NewChatMember)
	}

	if message.LeftChatMember != nil {
		err = r.removeUser(message.LeftChatMember.ID)
	}

	return err
}

func (r *Redis) saveUser(user *messages.User) error {
	r.log.Debugf("Saving user: %+v", *user)

	_, err := r.conn.Do("WATCH", r.userIdToString(user.ID))
	if err != nil {
		return errors.Wrapf(err, "redis command 'WATCH %s' failed", r.userIdToString(user.ID))
	}

	didUserChange, err := r.didUserChange(user)
	if !didUserChange {
		return nil
	}

	if err != nil && err != redis.ErrNil {
		return errors.Wrapf(err, "failed to detect if user '%+v' changed", user)
	}

	r.conn.Do("MULTI")
	defer r.conn.Do("EXEC")

	gob.Register(messages.User{})
	encodedUser, err := utils.EncodeObject(user)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize user '%+v'", user)
	}

	_, err = r.conn.Do("SET", r.userIdKey(user.ID), encodedUser)
	if err != nil {
		return errors.Wrapf(err, "redis command 'SET %s' failed", r.userIdKey(user.ID), encodedUser)
	}

	r.conn.Do("SADD", r.userFirstnameKey(user.FirstName), user.ID)
	if user.Username != "" { // Telegram username's are optional
		r.conn.Do("SADD", r.userUsernameKey(user.Username), user.ID)
	}

	return nil
}

func (r *Redis) didUserChange(user *messages.User) (bool, error) {
	user, err := r.getUser(user.ID)
	if err == redis.ErrNil {
		return true, redis.ErrNil
	}

	if err != nil {
		return true, errors.Wrapf(err, "failed to get user '%+v'", user)
	}

	return !reflect.DeepEqual(user, user), nil //TODO is == allowed for comparing structs, or reflect.DeepEqual()?
}

// Inverted Index
// chatid => [userid-DougC, userId-DougB, userId-Brian, userID-Tokie]
// firstname:doug => [userid-DougC, userid-DougB]
// return scores for both DougC and DougB since they are in => [userid-DougC, userid-DougB] and its an ambigous search term
//

// user:IsMemberOf:Room
// toDo: only remove user if its in none of the chat rooms...
func (r *Redis) removeUser(userId int64) error {
	_, err := r.conn.Do("WATCH", r.userIdToString(userId))
	if err != nil {
		return errors.Wrapf(err, "redis command 'WATCH %s' failed", r.userIdToString(userId))
	}

	user, err := r.getUser(userId)
	if err != nil {
		return errors.Wrapf(err, "failed to get user with userId '%d'", userId)
	}

	r.log.Debugf("Removing user: %v", user)
	r.conn.Do("MULTI")
	defer r.conn.Do("EXEC") // EXEC may fail if user changes....so need to retry entire operation

	_, err = r.conn.Do("SREM", r.userFirstnameKey(user.FirstName))
	if err != nil {
		return errors.Wrapf(err, "redis command 'SREM %s' failed", r.userFirstnameKey(user.FirstName))
	}

	_, err = r.conn.Do("SREM", r.userUsernameKey(user.Username))
	if err != nil {
		return errors.Wrapf(err, "redis command 'SREM %s' failed", r.userUsernameKey(user.Username))
	}

	_, err = r.conn.Do("DEL", r.userIdToString(userId))
	if err != nil {
		return errors.Wrapf(err, "redis command 'DEL %s' failed", r.userIdToString(userId))
	}

	return nil
}

func (r *Redis) GetUsersFromUsername(username string) ([]messages.User, error) {
	//	// TODO: this is not unique across all of telegram. There can be multiple users with the same username/firstName, so we need a more unique way to store it.
	reply, err := redis.Int64s(r.conn.Do("SMEMBERS", r.userUsernameKey(username)))
	if err != nil {
		return nil, errors.Wrapf(err, "redis command 'SMEMBERS %s' failed", r.userUsernameKey(username))
	}

	var users []messages.User
	for _, userId := range reply {
		user, err := r.getUser(userId)
		if err == redis.ErrNil {
			continue
		} else if err != nil {
			return nil, errors.Wrapf(err, "failed to get user with userId '%d'", userId)
		}

		users = append(users, *user) // UGH why does it think usersReply is not []int64 and thus userId an int64
	}

	return users, nil
}

func (r *Redis) GetUsersFromFirstName(firstName string) ([]messages.User, error) {
	// TODO: this is not unique across all of telegram. There can be multiple users with the same username/firstName, so we need a more unique way to store it.
	reply, err := redis.Int64s(r.conn.Do("SMEMBERS", r.userFirstnameKey(firstName)))
	if err != nil {
		return nil, errors.Wrapf(err, "redis command 'SMEMBERS %s' failed", r.userFirstnameKey(firstName))
	}

	var users []messages.User
	for _, userId := range reply {
		user, err := r.getUser(userId)
		if err == redis.ErrNil {
			continue
		} else if err != nil {
			return nil, errors.Wrapf(err, "failed to get user with userId '%d'", userId)
		}

		users = append(users, *user) // UGH why does it think usersReply is not []int64 and thus userId an int64
	}

	return users, nil
}

func (r *Redis) isBot(user *messages.User) (bool, error) {
	user, err := r.getUser(user.ID)
	if err != nil {
		return true, errors.Wrapf(err, "failed to get user with user '%+v'", user)
	}

	return user.IsBot && user == nil, nil
}

func (r *Redis) getUser(userId int64) (*messages.User, error) {
	reply, err := redis.Bytes(r.conn.Do("GET", r.userIdKey(userId))) // whats the diff between redis.Values and redis.StringMap in this case...?
	if reply == nil {
		return nil, redis.ErrNil
	}

	if err != nil {
		return nil, errors.Wrapf(err, "redis command 'GET %s' failed", r.userIdKey(userId))
	}

	// TODO: see what is stored if you try saving something that does not have a value and accessing it
	// Have it return in preference "firstname lastname", 2) "firstname" 3) username

	var user messages.User
	err = utils.DecodeObject(reply, &user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deserialize redis reply '%+v'", reply)
	}

	return &user, nil
}
