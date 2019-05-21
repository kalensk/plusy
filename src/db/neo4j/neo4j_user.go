package neo4j

import "github.com/kalensk/plusy/src/messages"

func (n *Neo4j) IncrementCount(chatId int64, dateGiven int64, giverUser *messages.User, receiverUser *messages.User) (string, error) {
	return "", nil
}

func (n *Neo4j) GetPointsReceived(chatId int64, giverUserId int64) ([]messages.UserPoints, error) {
	return []messages.UserPoints{}, nil
}

func (n *Neo4j) GetPointsGiven(chatId int64, receiverUserId int64) ([]messages.UserPoints, error) {
	return []messages.UserPoints{}, nil
}

func (n *Neo4j) GetTopN(chatId int64, num int) ([]messages.UserPoints, error) {
	return []messages.UserPoints{}, nil
}

func (n *Neo4j) SaveOrRemoveUser(message *messages.Message) error {
	return nil
}

func (n *Neo4j) GetUsersFromUsername(username string) ([]messages.User, error) {
	return []messages.User{}, nil
}

func (n *Neo4j) GetUsersFromFirstName(firstName string) ([]messages.User, error) {
	return []messages.User{}, nil
}
