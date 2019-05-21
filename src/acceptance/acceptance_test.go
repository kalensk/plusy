package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kalensk/plusy/src/db/redisdb"
	"github.com/kalensk/plusy/src/messages"
	"github.com/kalensk/plusy/src/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var docker *testutils.Docker
var database *redisdb.Redis
var log *logrus.Logger

func TestMain(m *testing.M) {
	docker = testutils.NewRedisDocker(m)
	log = docker.Log
	database = docker.SetupTestRedis()
	returnCode := docker.Run()
	docker.TearDown(returnCode)
}

func getFakeResponseBytes(fileName string) []byte {
	_, acceptanceTestFileName, _, _ := runtime.Caller(0)
	responseFile := filepath.Join(filepath.Dir(acceptanceTestFileName), fileName)
	fakeResponseBytes, err := ioutil.ReadFile(responseFile)
	if err != nil {
		panic(err)
	}

	//var update messages.UpdateResponse
	//err = json.Unmarshal(fakeResponseBytes, &update)

	return fakeResponseBytes
}

func Test_Acceptance(t *testing.T) {
	var command *exec.Cmd

	fakeResponse1 := getFakeResponseBytes("response1.json")
	//fakeResponse2 := getFakeResponseBytes("response2.json")

	var count int
	testServer := httptest.NewServer(http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		fmt.Println("Inside testServer")
		fmt.Println(request.Method)
		fmt.Println(request.RequestURI)
		fmt.Println(request.Body)

		if strings.Contains(request.URL.Path, "getUpdates") {
			respWriter.WriteHeader(200)
			respWriter.Header().Add("Content-Type", "application/json")
			respWriter.Write(fakeResponse1)
			if count > 0 {
				command.Process.Kill()
			}
			count++
		}

		if strings.Contains(request.URL.Path, "sendMessage") {
			requestBytes, err := ioutil.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(requestBytes))
		}

	}))
	defer testServer.Close()

	plusyBinary := filepath.Join(os.Getenv("GOPATH"), "bin", "plusy")
	command = exec.Command(plusyBinary, "--foreground",
		"--token", "123",
		"--telegram-api-url", testServer.URL+"/bot",
		"--database-host-port", docker.GetDatabaseConnectionString())

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	command.Stdout = &stdOut
	command.Stderr = &stdErr
	err := command.Start()
	duration, err := time.ParseDuration("10m")
	time.Sleep(duration)
	if err != nil {
		log.Println(err)
		log.Println(stdOut.String())
		log.Println(stdErr.String())
	}
}

func Test_AnotherAcceptanceTest(t *testing.T) {
	//plusOneAckMessage := messages.SendMessageRequest{ChatId: -276219865, Text: "fuzzie now has 41 plusies!"}
	//plusOneAckMessageBytes, err := json.Marshal(plusOneAckMessage)

	sendMessageResponse := messages.SendMessageResponse{Ok: true}
	sendMessageResponseBytes, err := json.Marshal(sendMessageResponse)

	fakeTelegramServer := testutils.NewFakeTelgramServer2(t)
	fakeTelegramServer.AddHandler("/getUpdates", func(respWriter http.ResponseWriter, request *http.Request) {
		respWriter.Write(getFakeResponseBytes("response1.json"))
	})
	fakeTelegramServer.AddHandler("/getUpdates", func(respWriter http.ResponseWriter, request *http.Request) {
		respWriter.Write(getFakeResponseBytes("response2.json"))
	})
	fakeTelegramServer.AddHandler("/sendMessage", func(respWriter http.ResponseWriter, request *http.Request) {
		respWriter.Write(sendMessageResponseBytes)
		assert.Equal(t, "sendMessageTest", request.URL.Path, "expected request path to match")
	})

	fakeTelegramServer.Start()
	defer fakeTelegramServer.Close()

	plusyBinary := filepath.Join(os.Getenv("GOPATH"), "bin", "plusy")
	command := exec.Command(plusyBinary,
		"--foreground",
		"--token", "123",
		"--telegram-api-url", fakeTelegramServer.Url()+"/bot",
		"--database-host-port", docker.GetDatabaseConnectionString())
	// todo: pass in a flag (-f) to write to stdout/stderr instead of plusy.log - DONE!

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	command.Stdout = &stdOut
	command.Stderr = &stdErr
	err = command.Start()
	duration, err := time.ParseDuration("2s")
	time.Sleep(duration)
	//if err != nil {
	log.Println(err)
	log.Println(stdOut.String())
	log.Println(stdErr.String())
	//}
}
