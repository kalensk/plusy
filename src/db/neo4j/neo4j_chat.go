package neo4j

import "github.com/kalensk/plusy/src/messages"

func (n *Neo4j) SaveLastMessage(message *messages.Message) error {
	return nil
}

func (n *Neo4j) GetLastMessage(chatId int64) (*messages.Message, error) {
	return &messages.Message{}, nil
}
