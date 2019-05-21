package client

import "github.com/kalensk/plusy/src/messages"

type ChatClient interface {
	GetUpdate(offset int64) (messages.UpdateResponse, error)
	GetUpdates(timeout int, offset int64) (messages.UpdateResponse, error)
	SendPlusOneAckMessage(chatId int64, userFirstName string, currentCount string) error
	ProcessCommand(message *messages.Message, cmdPosition int) error
}
