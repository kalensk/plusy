package db

import (
	"github.com/kalensk/plusy/src/messages"
)

type Database interface {
	IncrementCount(chatId int64, dateGiven int64, giverUser *messages.User, receiverUser *messages.User) (string, error)
	GetPointsReceived(FchatId int64, giverUserId int64) ([]messages.UserPoints, error)
	GetPointsGiven(chatId int64, receiverUserId int64) ([]messages.UserPoints, error)
	GetTopN(chatId int64, n int) ([]messages.UserPoints, error)

	SaveOffset(offset int64) error
	GetNextOffset() (int64, error)

	SaveLastMessage(message *messages.Message) error
	GetLastMessage(chatId int64) (*messages.Message, error)

	SaveOrRemoveUser(message *messages.Message) error
	GetUsersFromUsername(username string) ([]messages.User, error)
	GetUsersFromFirstName(firstName string) ([]messages.User, error)

	Ping() error
	Close() error
	ClearData() (string, error)
}
