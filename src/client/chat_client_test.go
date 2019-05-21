package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kalensk/plusy/src/db/redisdb"
	"github.com/kalensk/plusy/src/messages"
	"github.com/kalensk/plusy/src/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

/*
A new bot does not see the audio from another bot or something?

[TRACE] Processing flipTable command
[TRACE] Sending Telegram sendAudio request with message: {"chat_id":-276219865,"text":"","sticker":"","parse_mode":"","audio":"CQADAQADXwADhgbRRLdgKkZos0xuAg"}
[TRACE] Telegram response from sendAudio: {"ok":false,"error_code":400,"description":"Bad Request: wrong file identifier/HTTP URL specified"}
[TRACE] saving telemgram offset 28838563
[TRACE] returned telegram offset 28838563
*/

var database *redisdb.Redis
var log *logrus.Logger

func TestMain(m *testing.M) {
	docker := testutils.NewRedisDocker(m)
	log = docker.Log
	database = docker.SetupTestRedis()

	returnCode := docker.Run()
	docker.TearDown(returnCode)
}

var message = messages.Message{
	MessageID: 0,
	Chat:      testutils.SekretChatRoom,
	From:      testutils.UserDougCat,
	Text:      testutils.StringPointer("some message"),
}

func TestChatClient_SendPlusOneAckMessage(t *testing.T) {
	var actualPlusOneAckMessage messages.SendMessageRequest
	testServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(responseWriter, "Some Telegram Server Response")

		actualRequestBodyBytes, err := ioutil.ReadAll(request.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(actualRequestBodyBytes, &actualPlusOneAckMessage)
		assert.NoError(t, err)
	}))
	defer testServer.Close()

	client := New(log, testServer.URL+"/bot", "123456", database)
	client.SendPlusOneAckMessage(message.Chat.ID, message.From.FirstName, "2")

	expectedPlusOneAckMessage := messages.SendMessageRequest{
		ChatId: testutils.SekretChatRoom.ID,
		Text:   "doug has 2 plusies!!!"}
	assert.Equal(t, expectedPlusOneAckMessage, actualPlusOneAckMessage)
}
