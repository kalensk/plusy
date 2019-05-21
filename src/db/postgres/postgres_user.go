package postgres

import "github.com/kalensk/plusy/src/messages"

func (p *Postgres) IncrementCount(chatId int64, dateGiven int64, giverUser *messages.User, receiverUser *messages.User) (string, error) {
	return "", nil
}

func (p *Postgres) GetPointsReceived(chatId int64, giverUserId int64) ([]messages.UserPoints, error) {
	return []messages.UserPoints{}, nil
}

func (p *Postgres) GetPointsGiven(chatId int64, receiverUserId int64) ([]messages.UserPoints, error) {
	return []messages.UserPoints{}, nil
}

func (p *Postgres) GetTopN(chatId int64, n int) ([]messages.UserPoints, error) {
	return []messages.UserPoints{}, nil
}

func (p *Postgres) SaveOrRemoveUser(message *messages.Message) error {
	return nil
}

func (p *Postgres) GetUsersFromUsername(username string) ([]messages.User, error) {
	return []messages.User{}, nil
}

func (p *Postgres) GetUsersFromFirstName(firstName string) ([]messages.User, error) {
	return []messages.User{}, nil
}
