package plusy

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kalensk/plusy/src/client"
	"github.com/kalensk/plusy/src/db/redisdb"
	"github.com/kalensk/plusy/src/messages"
	"github.com/kalensk/plusy/src/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var database *redisdb.Redis
var log *logrus.Logger

func TestMain(m *testing.M) {
	docker := testutils.NewRedisDocker(m)
	log = docker.Log
	database = docker.SetupTestRedis()
	returnCode := docker.Run()
	docker.TearDown(returnCode)
}

var previousMessage = messages.Message{
	MessageID: 0,
	Chat:      testutils.SekretChatRoom,
	From:      testutils.UserDougCat,
	Text:      testutils.StringPointer("im a previous message"),
}

var plusOneReplyMessage = messages.Message{
	MessageID:      1,
	Chat:           testutils.SekretChatRoom,
	From:           testutils.UserDougBat,
	ReplyToMessage: &messages.ReplyToMessage{previousMessage},
	Text:           testutils.StringPointer("+1"),
}

var plusOneMessage = messages.Message{
	MessageID: 2,
	Chat:      testutils.SekretChatRoom,
	From:      testutils.UserMrRogers,
	Text:      testutils.StringPointer("+1"),
}

func Test_InlinePlusOne(t *testing.T) {
	database.ClearData()

	MockUpdateResponse := messages.UpdateResponse{
		Ok: true,
		Result: []messages.Result{
			{UpdateID: 0, Message: &previousMessage},
			{UpdateID: 1, Message: &plusOneMessage},
		},
	}

	mockClient := new(client.MockChatClient)
	mockClient.On("SendPlusOneAckMessage", mock.Anything, mock.Anything, mock.Anything)

	plusyApp := New(log, database, mockClient)
	plusyApp.ProcessResults(MockUpdateResponse.Result)

	chatID := previousMessage.Chat.ID

	actualPointsReceived, err := database.GetPointsReceived(chatID, previousMessage.From.ID)
	assert.NoError(t, err, "GetPointsReceived() succeeds on a reply message")
	assert.Equal(t, previousMessage.From, actualPointsReceived[0].User, "Getting user from GetPointsReceived() on reply message should succeed")
	assert.Equal(t, "1", actualPointsReceived[0].Plusies, "Getting plusies from GetPointsReceived() on reply message should succeed")
	mockClient.AssertExpectations(t)

	actualPointsGiven, err := database.GetPointsGiven(chatID, plusOneMessage.From.ID)
	assert.NoError(t, err, "GetPointsGiven() succeeds on a reply message")
	assert.Equal(t, plusOneMessage.From, actualPointsGiven[0].User, "Getting user from GetPointsGiven() on reply message should succeed")
	assert.Equal(t, "1", actualPointsGiven[0].Plusies, "Getting plusies from GetPointsGiven() on reply message should succeed")
	mockClient.AssertExpectations(t)
}

func Test_ReplyMessagePlusOne(t *testing.T) {
	database.ClearData()

	MockUpdateResponse := messages.UpdateResponse{
		Ok: true,
		Result: []messages.Result{
			{UpdateID: 0, Message: &previousMessage},
			{UpdateID: 1, Message: &plusOneReplyMessage},
		},
	}

	mockClient := new(client.MockChatClient)
	mockClient.On("SendPlusOneAckMessage", mock.Anything, mock.Anything, mock.Anything)

	plusyApp := New(log, database, mockClient)
	plusyApp.ProcessResults(MockUpdateResponse.Result)

	chatID := previousMessage.Chat.ID

	actualPointsReceived, err := database.GetPointsReceived(chatID, previousMessage.From.ID)
	assert.NoError(t, err, "GetPointsReceived() succeeds on a reply message")
	assert.Equal(t, previousMessage.From, actualPointsReceived[0].User, "Getting user from GetPointsReceived() on reply message should succeed")
	assert.Equal(t, "1", actualPointsReceived[0].Plusies, "Getting plusies from GetPointsReceived() on reply message should succeed")
	mockClient.AssertExpectations(t)

	actualPointsGiven, err := database.GetPointsGiven(chatID, plusOneReplyMessage.From.ID)
	assert.NoError(t, err, "GetPointsGiven() succeeds on a reply message")
	assert.Equal(t, plusOneReplyMessage.From, actualPointsGiven[0].User, "Getting user from GetPointsGiven() on reply message should succeed")
	assert.Equal(t, "1", actualPointsGiven[0].Plusies, "Getting plusies from GetPointsGiven() on reply message should succeed")
	mockClient.AssertExpectations(t)
}

func Test_SaveLastMessageSucceedsForEditedMessage(t *testing.T) {
	database.ClearData()

	editedMessage := &messages.Message{
		MessageID: 0,
		Chat:      testutils.SekretChatRoom,
		EditDate:  testutils.Int64Pointer(1537927460),
		From:      testutils.UserDougCat,
		Text:      testutils.StringPointer("i am an edited message"),
	}

	MockUpdateResponse := messages.UpdateResponse{
		Ok: true,
		Result: []messages.Result{
			{UpdateID: 0, EditedMessage: editedMessage},
		},
	}

	plusyApp := New(log, database, new(client.MockChatClient))
	plusyApp.ProcessResults(MockUpdateResponse.Result)

	actualLastMessage, _ := database.GetLastMessage(testutils.SekretChatRoom.ID)
	assert.Equal(t, editedMessage, actualLastMessage, "Saving the last message succeeds when the last message returned from getUpdates is an edited message")
}

func Test_getPreviousUser_SucceedsWhenThereIsMoreThanOneMessage(t *testing.T) {
	database.ClearData()
	plusyApp := New(log, database, new(client.MockChatClient))

	results := []messages.Result{
		{UpdateID: 0, Message: &previousMessage},
		{UpdateID: 1, Message: &plusOneMessage},
	}

	previousUser, _ := plusyApp.getPreviousUser(1, results, results[1].Message)
	assert.Equal(t, testutils.UserDougCat, previousUser)
}

func Test_getPreviousUser_SucceedsWhenThereIsOneMessage(t *testing.T) {
	database.ClearData()
	database.SaveLastMessage(&plusOneReplyMessage)

	plusyApp := New(log, database, new(client.MockChatClient))

	results := []messages.Result{
		{UpdateID: 0, Message: &previousMessage},
	}

	previousUser, _ := plusyApp.getPreviousUser(0, results, results[0].Message)
	assert.Equal(t, testutils.UserDougBat, previousUser)
}

func Test_getPreviousUser_ReturnsAnErrorWhenThereIsOneMessageAndTheDatabaseFails(t *testing.T) {
	database.ClearData()
	plusyApp := New(log, database, new(client.MockChatClient))

	results := []messages.Result{
		{UpdateID: 0, Message: &previousMessage},
	}

	previousUser, err := plusyApp.getPreviousUser(0, results, results[0].Message)
	assert.EqualError(t, err, "failed to get last message for chat room: 789: redigo: nil returned")
	assert.Nil(t, previousUser)
}

func Test_ExponentialRetry(t *testing.T) {
	t.SkipNow()
	database.ClearData()

	duration, err := time.ParseDuration("500ms")
	if err != nil {
		panic(err)
	}


	// f(x) = 5*x/(x + 5)
	// 2^f(x)

	// sleep(2^(x*x/(x*x+50)))

	//sum := 1
	max_retry := 1.0
	for true {
		for j := 1.0; j <= max_retry; j++ {
			//sum = sum * 2
			sum := math.Pow(2.0, max_retry)
			fmt.Printf("Message %f : %f\n", max_retry, sum)
			max_retry++
			time.Sleep(duration)
		}
	}


}
