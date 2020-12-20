package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/kalensk/plusy/src/db"
	"github.com/kalensk/plusy/src/messages"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Telegram struct {
	log        *logrus.Logger
	url        string
	database   db.Database // why can't I make this *db.Database
	stickers   []string
	soundBites []string
	quotes     []string
}

func New(log *logrus.Logger, apiUrl string, token string, database db.Database) *Telegram {
	url := apiUrl + token + "/"

	// ToDo: populate with a call to
	// GET 'https://api.telegram.org/<api-token>/getStickerSet?name=DonutAndCoffee'
	//  jq '.result.stickers | .[] | .file_id'

	// DonutAndCoffee sticker pack
	stickers := []string{
		"CAADAgADhgIAAkcVaAmAJ7NbemcocAI",
		"CAADAgADtQIAAkcVaAk4LbxZ_Ei6mQI",
		"CAADAgADiAIAAkcVaAkNL7KdtavmWAI",
		"CAADAgADkgIAAkcVaAlpiUuV3-CK5gI",
		"CAADAgADigIAAkcVaAkO1Lftas8KSQI",
		"CAADAgADiwIAAkcVaAl6-6cfPdXGmgI",
		"CAADAgADmAIAAkcVaAl4gSxYpUlsFAI",
		"CAADAgADnAIAAkcVaAlpX0t5hd0HaAI",
		"CAADAgADrwIAAkcVaAkwqZ0sPGF-3wI",
		"CAADAgADjgIAAkcVaAnKte4F9usSIgI",
		"CAADAgADlAIAAkcVaAmR876wugABwasC",
		"CAADAgADsAIAAkcVaAmXeSFxP-EqDwI",
		"CAADAgADmgIAAkcVaAkaFwoozllspQI",
		"CAADAgADjwIAAkcVaAk_f_rRm4A5bQI",
		"CAADAgADkAIAAkcVaAn7uW7oNOJ9YAI",
		"CAADAgADlwIAAkcVaAmH1nLynvme6AI",
		"CAADAgADkQIAAkcVaAnWwvZMwGxOcwI",
		"CAADAgADlQIAAkcVaAkG5pdrTGfBvAI",
		"CAADAgADlgIAAkcVaAnwuBPg4BZQmgI",
		"CAADAgADmQIAAkcVaAm8HwZufZummgI",
		"CAADAgADmwIAAkcVaAlHvbtAMfqtawI",
		"CAADAgADjQIAAkcVaAn7gTOMFUhoewI",
		"CAADAgADnQIAAkcVaAkX4M7yYGW8MAI",
		"CAADAgADngIAAkcVaAkxPVk2Gt0t3AI",
		"CAADAgADnwIAAkcVaAnW2tGaOkkU7wI",
		"CAADAgADoAIAAkcVaAmG5GLznokYGgI",
		"CAADAgADoQIAAkcVaAkWoCnNtKe5dAI",
		"CAADAgADogIAAkcVaAm_l8R0g8wieQI",
		"CAADAgADpAIAAkcVaAlR0TGwaHTK_QI",
		"CAADAgADpQIAAkcVaAn8a52QxbSpfQI",
		"CAADAgADpgIAAkcVaAkpenv_Dhq_JAI",
		"CAADAgADpwIAAkcVaAlBFKBsTr43JgI",
		"CAADAgADqAIAAkcVaAllGKpjZDCAgQI",
		"CAADAgADqQIAAkcVaAlHBy_HsaF8HgI",
		"CAADAgADqgIAAkcVaAkAAbFNXxuKRv4C",
		"CAADAgADrAIAAkcVaAn8e3QPPV8zAgI",
		"CAADAgADrQIAAkcVaAn1loAVfPjt9QI",
		"CAADAgADrgIAAkcVaAmFU5Elxy-lVgI",
		"CAADAgADsQIAAkcVaAkh7PN7i7ZOSwI",
		"CAADAgADsgIAAkcVaAnzN6EnIlkusAI",
	}

	soundBites := []string{
		"CQADAQADWAADhgbRRKLr1SLyuyA8Ag", // blup
		"CQADAQADXwADhgbRRLdgKkZos0xuAg", // ah
		"CQADAQADYAADhgbRRJmJPUjKxP-bAg", // od2
		"CQADAQADYgADhgbRRF4LYoOtKBxJAg", // hackedb
		"CQADAQADYwADhgbRRONLrzcJXDFzAg", // spin
		"CQADAQADZAADhgbRREu8rnXbgFCIAg", // fired
		"CQADBAADPgADtNPUULQxO7DhGCRsAg", // cookies
	}

	quotes := []string{
		`When in doubt, use brute force.
         Ken Thompson
         Bell Labs`,

		`Allocate four digits for the year part of a date: a new
millenium is coming.
        David Martin
        Norristown, Petmsylvania`,
	}

	return &Telegram{log: log, url: url, database: database, stickers: stickers, soundBites: soundBites, quotes: quotes}
}

func (t *Telegram) GetUpdate(offset int64) (messages.UpdateResponse, error) {
	return t.getUpdates(0, offset, 1)
}

func (t *Telegram) GetUpdates(timeout int, offset int64) (messages.UpdateResponse, error) {
	return t.getUpdates(timeout, offset, -1)
}

func (t *Telegram) getUpdates(timeout int, offset int64, limit int) (messages.UpdateResponse, error) {
	//params := fmt.Sprintf("?allowed_messages=[\"message\"]&timeout=%d&offset=%d", timeout, offset)
	params := fmt.Sprintf("?timeout=%d&offset=%d", timeout, offset)
	if limit > 0 {
		params = params + "&limit=" + string(limit) // why not `params += "&limit=" + string(limit)` instead?
	}

	resp, err := http.Get(t.url + "getUpdates" + params)
	if err != nil {
		return messages.UpdateResponse{}, errors.Wrapf(err, "failed to get telegram updates to /%s", "getUpdates"+params) // don't log the url which has the api token
		// ToDo: exponential retry while logging once an Error. Add as a helper function so http.Post can use it.
	}

	/*
				if internet connect goes down then DNS lookup for api.telegram.org fails and you didn't handle the error

				panic: Get https://api.telegram.org/...: dial tcp: lookup api.telegram.org on 127.0.1.1:53: read udp 127.0.0.1:46463->127.0.1.1:53: i/o timeout

		goroutine 1 [running]:
		github.com/kalensk/plusy/src/client.(*Telegram).getUpdates(0xc4200c9180, 0x2, 0x1b80f84, 0xffffffffffffffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
			/home/kalensk/go/src/github.com/kalensk/plusy/src/client/telegram.go:118 +0x482
		github.com/kalensk/plusy/src/client.(*Telegram).GetUpdates(0xc4200c9180, 0x2, 0x1b80f84, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
			/home/kalensk/go/src/github.com/kalensk/plusy/src/client/telegram.go:106 +0x8a
		main.main()
			/home/kalensk/go/src/github.com/kalensk/plusy/src/main.go:115 +0x1b2
		exit status 2
	*/

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return messages.UpdateResponse{}, errors.Wrap(err, "failed to read response body from call to getUpdates")
	}

	var update messages.UpdateResponse
	err = json.Unmarshal(respBytes, &update)
	if err != nil {
		return messages.UpdateResponse{}, errors.Wrap(err, "failed to deserialize response bytes from call to getUpdates")
	}

	if len(update.Result) > 0 {
		t.log.Debugf("Getting Telegram updates: %s", string(respBytes))
	}

	return update, nil
}

func (t *Telegram) SendPlusOneAckMessage(chatId int64, userFirstName string, currentCount string) error {
	text := fmt.Sprintf("%s has %s plusies!", userFirstName, currentCount)
	if currentCount == "1" {
		text = fmt.Sprintf("%s has %s plusy!", userFirstName, currentCount)
	}

	plusOneAckMessage := messages.SendMessageRequest{ChatId: chatId, Text: text}
	plusOneAckMessageBytes, err := json.Marshal(plusOneAckMessage)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize message %+v", plusOneAckMessage)
	}

	return t.sendMessage(plusOneAckMessageBytes)
}

func (t *Telegram) sendMessage(message []byte) error {
	return t.send("sendMessage", message)
}

func (t *Telegram) sendSticker(message []byte) error {
	return t.send("sendSticker", message)
}

func (t *Telegram) sendAudio(message []byte) error {
	return t.send("sendAudio", message)
}

func (t *Telegram) send(method string, message []byte) error {
	t.log.Debugf("Sending Telegram %s request with message: %s\n", method, string(message))
	resp, err := http.Post(t.url+method, "application/json", bytes.NewReader(message))
	if err != nil {
		return errors.Wrapf(err, "failed to POST message: %+v to: %s", message, t.url+method)
		// ToDo: exponential retry while logging once an Error.  Add as a helper function so http.Get can use it.
		// https://cloud.google.com/iot/docs/how-tos/exponential-backoff
		// https://docs.microsoft.com/en-us/dotnet/standard/microservices-architecture/implement-resilient-applications/explore-custom-http-call-retries-exponential-backoff
		// https://developers.google.com/api-client-library/java/google-http-java-client/reference/1.20.0/com/google/api/client/util/ExponentialBackOff
		// https://github.com/googleapis/google-http-java-client/blob/master/google-http-client/src/main/java/com/google/api/client/util/ExponentialBackOff.java
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	t.log.Debugf("Telegram response from %s: %s\n", method, string(respBytes))
	return nil
}

func (t *Telegram) ProcessCommand(message *messages.Message, cmdPosition int) error {
	t.log.Debugf("Received command: %s\n", *message.Text)

	chatId := message.Chat.ID

	text := *message.Text
	textRune := []rune(text)
	command := string(textRune[0:cmdPosition]) // ToDo: write test first and then change this to 1:cmdPosition and then remove the leading / in the switch statement
	command = strings.SplitN(command, "@", 2)[0] // normalize command since it can be "/stats@plusyTestBot"

	args := strings.TrimSpace(string(textRune[cmdPosition:]))

	switch command {
	case "/stats":
		return t.processStatsCommand(chatId, args)
	case "/topstats", "/top":
		return t.processTopStatsCommand(chatId)
	case "/donut":
		return t.processDonutCommand(chatId)
	case "/fliptable", "/flipstable", "/sound":
		return t.processSoundCommand(chatId)
	case "/quote":
		return t.processQuoteCommand(chatId)
	case "/help":
		return t.processHelpCommand(chatId)
	case "/halp":
		// send cat hanging pic
	}

	return nil
}

// ToDo: sort results before returning
func (t *Telegram) processTopStatsCommand(chatId int64) error {
	t.log.Debugf("Processing top stats command for chatId: %d\n", chatId)
	topN, err := t.database.GetTopN(chatId, 5)
	if err != nil {
		return errors.Wrapf(err, "failed to get top scores for chat %s (%d)", "friendlyChatName", chatId)
	}

	var text strings.Builder
	text.WriteString("Top Stats:\n")
	for _, pointStruct := range topN {
		if pointStruct.Plusies != "1" {
			text.WriteString(
				fmt.Sprintf("    %s has %s plusies\n", pointStruct.User.FirstName, pointStruct.Plusies))
		} else {
			text.WriteString(
				fmt.Sprintf("    %s has %s plusy\n", pointStruct.User.FirstName, pointStruct.Plusies))
		}
	}

	topNStatsRequest := messages.SendMessageRequest{ChatId: chatId, Text: text.String()}
	topNStatsRequestBytes, err := json.Marshal(topNStatsRequest)
	if err != nil {
		return errors.Wrap(err, "failed to serialize topN stats request")
	}

	t.sendMessage(topNStatsRequestBytes)
	return nil
}

func (t *Telegram) processStatsCommand(chatId int64, queriedUser string) error {
	t.log.Debugf("Processing stats command for queried user %s in chatId: %d\n", queriedUser, chatId)
	if queriedUser == "" {
		t.processTopStatsCommand(chatId)
		return nil
	}

	var users []messages.User
	if strings.HasPrefix(queriedUser, "@") {
		usersFromUserName, err := t.database.GetUsersFromUsername(queriedUser)
		if err != nil {
			return errors.Wrapf(err, "failed to get users with username '%s'", queriedUser)
		}

		users = usersFromUserName
	} else {
		usersFromFirstName, err := t.database.GetUsersFromFirstName(queriedUser)
		if err != nil {
			return errors.Wrapf(err, "failed to get users with firstname '%s'", queriedUser)
		}

		users = usersFromFirstName
	}

	t.log.Debugf("Processing stats command for users: %v", users)
	// ToDo: Return break down of points received by and given

	// /stats tokie
	// tokie has 2 plusies

	// Received By:
	//   doug 1
	//   brian 1
	//
	// Given 1 plusy:
	//   kaboodle 1

	var sb strings.Builder
	for _, user := range users {
		// write points received
		userPointsReceived, err := t.database.GetPointsReceived(chatId, user.ID)
		if err != nil {
			return errors.Wrapf(err, "failed to get points received for user %+v with chatId %d", user, chatId)
		}

		totalPlusies, err := t.sumPoints(userPointsReceived)
		if err != nil {
			return errors.Wrapf(err, "failed to sum points for user %+v", user)
		}

		sb.WriteString(fmt.Sprintf("%s has %s plusies\n", user.FirstName, totalPlusies))
		sb.WriteString("    Received by:\n")

		for _, userPoint := range userPointsReceived {
			sb.WriteString(fmt.Sprintf("        %s %s\n", userPoint.User.FirstName, userPoint.Plusies))
		}

		sb.WriteString("    Given to:\n")
		userPointsGiven, err := t.database.GetPointsGiven(chatId, user.ID)
		if err != nil {
			return errors.Wrapf(err, "failed to get points given for user %+v with chatId %d", user, chatId)
		}
		for _, userPoint := range userPointsGiven {
			sb.WriteString(fmt.Sprintf("        %s %s\n", userPoint.User.FirstName, userPoint.Plusies))
		}
	}

	statsRequest := messages.SendMessageRequest{ChatId: chatId, Text: sb.String()}
	statsRequestBytes, err := json.Marshal(statsRequest)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize stats request %+v", statsRequest)
	}

	t.sendMessage(statsRequestBytes)
	return nil
}

func (t *Telegram) sumPoints(userPoints []messages.UserPoints) (string, error) {
	var plusies int64
	for _, point := range userPoints {
		plusy, err := strconv.ParseInt(point.Plusies, 10, 64)
		if err != nil {
			return "", errors.Wrapf(err, "Failed to string plusy into int64 plusy: %v", point)
		}

		plusies += plusy
	}

	return strconv.FormatInt(plusies, 10), nil
}

func (t *Telegram) processDonutCommand(chatId int64) error {
	t.log.Debugf("Processing donut command\n")
	// ToDo: return random donut sticker from pack. use sticker.file_id and not sticker.thumb.file_id
	sticker := t.stickers[rand.Int()%len(t.stickers)]
	donutRequest := messages.SendMessageRequest{ChatId: chatId, Sticker: sticker} // "set_name": "DonutAndCoffee",
	donutRequestBytes, err := json.Marshal(donutRequest)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize donut request %+v", donutRequest)
	}

	t.sendSticker(donutRequestBytes)
	return nil
}

func (t *Telegram) processSoundCommand(chatId int64) error {
	t.log.Debugf("Processing flipTable command\n")
	soundBite := t.soundBites[rand.Int()%len(t.soundBites)]
	audioRequest := messages.SendMessageRequest{ChatId: chatId, Audio: soundBite}
	audioRequestBytes, err := json.Marshal(audioRequest)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize sound request %+v", audioRequest)
	}

	t.sendAudio(audioRequestBytes)
	return nil
}

func (t *Telegram) processQuoteCommand(chatId int64) error {
	t.log.Debugf("Processing quote command\n")
	quote := t.quotes[rand.Int()%len(t.quotes)]
	quoteRequest := messages.SendMessageRequest{ChatId: chatId, Text: quote}
	quoteRequestBytes, err := json.Marshal(quoteRequest)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize quote request %+v", quoteRequest)
	}

	t.sendMessage(quoteRequestBytes)
	return nil
}

func (t *Telegram) processHelpCommand(chatId int64) error {
	t.log.Debugf("Processing help command\n")
	text := `/stats [firstname | username]
/stats
/topstats
/quote
/donut
/fliptable | /flipstable | /sound
/halp
/help
`
	helpRequest := messages.SendMessageRequest{ChatId: chatId, Text: text}
	helpRequestBytes, err := json.Marshal(helpRequest)
	if err != nil {
		return errors.Wrapf(err, "failed to serialize help request %+v", helpRequest)
	}

	t.sendMessage(helpRequestBytes)
	return nil
}
