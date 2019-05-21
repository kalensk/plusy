package db

import (
	"reflect"
	"testing"

	"github.com/kalensk/plusy/src/db/redisdb"
	"github.com/kalensk/plusy/src/messages"
	"github.com/kalensk/plusy/src/testutils"
	"github.com/stretchr/testify/assert"
)

var database *redisdb.Redis

func TestMain(m *testing.M) {
	docker := testutils.NewRedisDocker(m)
	database = docker.SetupTestRedis()
	returnCode := docker.Run()
	docker.TearDown(returnCode)
}

func TestDatabase_GetUsersFromFirstName(t *testing.T) {
	dougCatMessage := &messages.Message{
		MessageID: 0,
		Chat:      testutils.SekretChatRoom,
		From:      testutils.UserDougCat,
		Text:      testutils.StringPointer("i am a message from doug_cat"),
	}

	dougBatMessage := &messages.Message{
		MessageID: 0,
		Chat:      testutils.SekretChatRoom,
		From:      testutils.UserDougBat,
		Text:      testutils.StringPointer("i am a message from doug_bat"),
	}

	mrRogersMessaage := &messages.Message{
		MessageID: 0,
		Chat:      testutils.SekretChatRoom,
		From:      testutils.UserMrRogers,
		Text:      testutils.StringPointer("i am a message from mr_rogers"),
	}

	database.SaveOrRemoveUser(dougCatMessage)
	database.SaveOrRemoveUser(dougBatMessage)
	database.SaveOrRemoveUser(mrRogersMessaage)

	expectedUsers, _ := database.GetUsersFromFirstName("mr")
	actualUsers := []messages.User{*testutils.UserMrRogers}
	assert.Equal(t, expectedUsers, actualUsers, "Should be able to retrieve saved user")

	expectedUsers, _ = database.GetUsersFromFirstName("doug")
	actualUsers = []messages.User{*testutils.UserDougCat, *testutils.UserDougBat}
	//assert.Equal(t, expectedUsers, actualUsers, "Should return all users when searching by first name")
	if !reflect.DeepEqual(expectedUsers, actualUsers) {
		t.Errorf("Expected to return all users when seraching by first name. Expected %v to Equal %v", expectedUsers, actualUsers)
	}
}

//// Dont need this test since there is no type "EditedMessage"
//// See similar test in acceptance_test.go
//func TestDatabase_SaveLastMessageSucceedsForEditedMessage(t *testing.T) {
//	chatSekretRoom := messages.Chat{ID: 789, Title: "sekret-room"}
//
//	editedMessage := &messages.Message{
//		MessageID: 0,
//		Chat:       chatSekretRoom,
//		EditDate:  int64Pointer(1537927460),
//		From:      &messages.User{ID: 2, FirstName: "mr", Username: "mr_rogers"},
//		Text:      stringPointer("i am an edited message"),
//	}
//
//	database.SaveLastMessage(editedMessage)
//	actualLastMessage := database.GetLastMessage(chatSekretRoom.ID)
//	assert.Equal(t, editedMessage, actualLastMessage, "Saving last message should work for edited messages")
//}
