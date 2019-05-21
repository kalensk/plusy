package main

import (
	"github.com/kalensk/plusy/src/db"
	"github.com/kalensk/plusy/src/db/neo4j"
	"github.com/kalensk/plusy/src/db/postgres"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/kalensk/plusy/src/client"
	"github.com/kalensk/plusy/src/db/redisdb"
	"github.com/kalensk/plusy/src/options"
	"github.com/kalensk/plusy/src/plusy"
)

func main() {
	isConnected := true

	opts := options.Parse()
	log := opts.Logger
	//chatIdToMessages := make(map[int64][]messages.Result)
	log.Info("starting plusy")

	database, err := getDatabase(opts)
	if err != nil {
		// ToDo: exponential retry of database forever
		log.Fatal("failed to connect to database due to: ", err)
	}

	chatClient := client.New(log, opts.TelegramApiUrl, opts.TelegramBotToken, database)
	plusyApp := plusy.New(log, database, chatClient)

	for isConnected {
		nextOffset, err := database.GetNextOffset()
		if err != nil {
			log.Error("failed to get next offset due to: ", err)
		}

		updates, err := chatClient.GetUpdates(2, nextOffset)
		if err != nil {
			log.Error("failed to get updates due to: ", err)
			continue
		}

		if !updates.Ok {
			log.Errorf("%d failed to get telegram updates due to: %s\n", updates.ErrorCode, updates.ErrorDescription)
			continue
		}

		if len(updates.Result) < 1 { // return if there are no messages
			continue
		}

		plusyApp.ProcessResults(updates.Result)
		time.Sleep(2 * time.Second)

		//for _, result := range updates.Result {
		//	if results, ok := chatIdToMessages[result.Message.Chat.ID]; ok {
		//		chatIdToMessages[result.Message.Chat.ID] = append(results, result)
		//	} else {
		//		chatIdToMessages[result.Message.Chat.ID] = []messages.Result{result}
		//	}
		//}
		//
		//
		//for _, results := range chatIdToMessages {
		//	go ProcessResults(database, chatClient, results)
		//}
	}

	database.Close()
}

func getDatabase(opts *options.Options) (db.Database, error) {
	switch strings.ToLower(opts.Database) {
	case "redis":
		return redisdb.New(opts.Logger, opts.DatabaseAddrPort)
	case "postgres":
		return postgres.New(opts.Logger, opts.DatabaseAddrPort)
	case "neo4j":
		return neo4j.New(opts.Logger, opts.DatabaseAddrPort)
	}

	return nil, errors.Errorf("Expected one of the following databases to use: redis, postgres, or neo4j. Found '%s' instead.", opts.Database)
}
