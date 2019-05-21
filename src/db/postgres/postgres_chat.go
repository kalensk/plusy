package postgres

import "github.com/kalensk/plusy/src/messages"

func (p *Postgres) SaveLastMessage(message *messages.Message) error {
	return nil
}

func (p *Postgres) GetLastMessage(chatId int64) (*messages.Message, error) {
	return &messages.Message{}, nil
}
