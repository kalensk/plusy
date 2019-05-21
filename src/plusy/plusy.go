package plusy

import (
	"strconv"
	"time"

	"github.com/kalensk/plusy/src/client"
	"github.com/kalensk/plusy/src/db"
	"github.com/kalensk/plusy/src/messages"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Plusy struct {
	log      *logrus.Logger
	database db.Database
	client   client.ChatClient
}

func New(log *logrus.Logger, database db.Database, client client.ChatClient) *Plusy {
	return &Plusy{log: log, database: database, client: client}
}

// ToDO: this should not return an error, correct?
func (p *Plusy) ProcessResults(results []messages.Result) {
	var message *messages.Message

	for i := 0; i < len(results); i++ { // replace with `for i, result := range results`?
		message = results[i].Message
		if !messages.HasText(message) {
			message = results[i].EditedMessage
			continue
		}

		p.log.WithFields(logrus.Fields{
			"chat_title": message.Chat.Title,
			"chat_id":    strconv.FormatInt(message.Chat.ID, 10),
			"user_id":    strconv.FormatInt(message.From.ID, 10),
			"first_name": message.From.FirstName,
			"user_name":  message.From.Username,
			"text":       *message.Text,
			"date":       time.Unix(message.Date, 0),
		}).Debug("Received message")

		p.database.SaveOrRemoveUser(message)

		if *message.Text == "+1" {
			// user of previous message if inline, or user of RepliedMessage
			previousUser, err := p.getPreviousUser(i, results, message)
			if err != nil {
				p.log.WithFields(logrus.Fields{
					"chat_title":    message.Chat.Title,
					"chat_id":       strconv.FormatInt(message.Chat.ID, 10),
					"date":          time.Unix(message.Date, 0),
					"giver_user":    message.From,
					"receiver_user": previousUser,
					"text":          *message.Text,
				}).Errorf("Failed to get previous user due to: ", err.Error())
				continue
			}

			if p.isSelfPlusOne(message.From.ID, previousUser.ID) {
				continue // don't allow self +1's
			}

			currentCount, err := p.database.IncrementCount(message.Chat.ID, message.Date, message.From, previousUser)
			if err != nil {
				p.log.WithFields(logrus.Fields{
					"chat_title":    message.Chat.Title,
					"chat_id":       strconv.FormatInt(message.Chat.ID, 10),
					"date":          time.Unix(message.Date, 0),
					"giver_user":    message.From,
					"receiver_user": previousUser,
					"text":          *message.Text,
				}).Errorf("Failed to increment count due to: ", err.Error())
				continue
			}

			err = p.client.SendPlusOneAckMessage(message.Chat.ID, previousUser.FirstName, currentCount)
			if err != nil {
				p.log.Error("Failed to send plus one acknowledgment message due to: ", err)
				continue
			}
			continue
		}

		cmdPosition, err := p.getCommandPosition(message)
		if err != nil {
			continue
		}

		err = p.client.ProcessCommand(message, cmdPosition)
		if err != nil {
			p.log.Errorf("Failed to process command for message %v", message)
			continue
		}
	}

	if message != nil { // guard against NewChatMember and other updates
		p.database.SaveLastMessage(message)
		offset := results[len(results)-1].UpdateID
		err := p.database.SaveOffset(offset)
		if err != nil {
			p.log.Errorf("Failed to save offset %d due to: ", offset, err.Error())
		}
	}
}

func (p *Plusy) isAReplyMessage(message *messages.Message) bool {
	return message.ReplyToMessage != nil
}

func (p *Plusy) isAnInlineMessage(message *messages.Message) bool {
	return message.ReplyToMessage == nil
}

func (p *Plusy) getPreviousUser(messageIndex int, results []messages.Result, currentMessage *messages.Message) (*messages.User, error) {
	if p.isAReplyMessage(currentMessage) {
		return currentMessage.ReplyToMessage.From, nil
	}

	if p.isAnInlineMessage(currentMessage) && messageIndex == 0 { // if we only got one update message then ask for the last saved message
		previousMessage, err := p.database.GetLastMessage(currentMessage.Chat.ID)
		if err != nil {
			return nil, err
		}

		return previousMessage.From, nil
	}

	if p.isAnInlineMessage(currentMessage) {
		return results[messageIndex-1].Message.From, nil
	}

	return nil, errors.New("failed to get previous user")
}

func (p *Plusy) isSelfPlusOne(userId int64, userIdOfPreviousMessageOrRepliedMessage int64) bool {
	return userId == userIdOfPreviousMessageOrRepliedMessage
}

func (p *Plusy) getCommandPosition(message *messages.Message) (int, error) {
	for _, entity := range message.Entities {
		if entity.Type == "bot_command" {
			return entity.Length, nil
		}
	}

	return -1, errors.Errorf("plusy command not found in message: %v", message.Entities)
}
