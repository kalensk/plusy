package redisdb

import (
	"encoding/gob"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/kalensk/plusy/src/messages"
	"github.com/kalensk/plusy/src/utils"
	"github.com/pkg/errors"
)

func (r *Redis) SaveLastMessage(message *messages.Message) error {
	if message == nil {
		return nil
	}

	gob.Register(messages.Message{})
	encodedMessage, err := utils.EncodeObject(message) // is this is saving null fields for some reason... for example: "{\"message_id\":26178,\"from\":{\"id\":308188500,\"is_bot\":false,\"first_name\":\"brain\",\"last_name\":\"fister\",\"username\":\"\"},\"date\":1537588507,\"chat\":{\"id\":-1001082930701,\"type\":\"supergroup\",\"title\":\"unwhirled\",\"all_members_are_administrators\":false},\"text\":\"dawg\",\"entities\":null,\"reply_to_message\":null,\"new_chat_member\":null,\"left_chat_member\":null,\"edit_date\":0}"
	if err != nil {
		return errors.Wrapf(err, "failed to serialize message '%+v' ", message)
	}

	key := "chat:lastMsg:" + strconv.FormatInt(message.Chat.ID, 10)
	_, err = r.conn.Do("SET", key, encodedMessage)
	if err != nil {
		panic(err)
		return errors.Wrapf(err, "redis command 'SET %s %s' failed", key, encodedMessage)
	}

	return nil
}

func (r *Redis) GetLastMessage(chatId int64) (*messages.Message, error) {
	key := "chat:lastMsg:" + strconv.FormatInt(chatId, 10)
	reply, err := redis.Bytes(r.conn.Do("GET", key))
	if err != nil {
		return nil, errors.Wrapf(err, "redis command 'GET %s' failed", key)
	}

	if reply == nil {
		return nil, errors.Wrapf(err, "redis command 'GET %s %s' returned nil", key)
	}

	var lastMessage messages.Message
	err = utils.DecodeObject(reply, &lastMessage)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to deserialize redis reply '%+v'", reply)
	}

	return &lastMessage, nil
}
