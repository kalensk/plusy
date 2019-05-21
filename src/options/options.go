package options

import (
	"flag"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"

	"github.com/sirupsen/logrus"
)

type Options struct {
	Logger           *logrus.Logger
	TelegramBotToken string
	TelegramApiUrl   string
	Database         string
	DatabaseAddrPort string
	Foreground       bool
}

func Parse() *Options {
	opts := Options{}

	flag.StringVar(&opts.TelegramBotToken, "token", "", "required: telegram bot authorization token")
	flag.StringVar(&opts.TelegramApiUrl, "telegram-api-url", "https://api.telegram.org/bot", "telegram bot api url")
	flag.StringVar(&opts.Database, "database", "redis", "database to use either: redis, postgres, or neo4j")
	flag.StringVar(&opts.DatabaseAddrPort, "database-host-port", "localhost:6379", "database host and port")
	flag.BoolVar(&opts.Foreground, "foreground", false, "run in foreground. don't log.")
	flag.Parse()

	opts.Logger = initializeLogger(opts.Foreground)

	if opts.TelegramBotToken == "" {
		log.Errorf("Missing api token. See --help")
		os.Exit(64)
	}

	rand.Seed(time.Now().UnixNano()) // this is not an option, move it some place more relevant

	return &opts
}

func initializeLogger(foreground bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter.(*logrus.TextFormatter).DisableLevelTruncation = false

	if foreground {
		logger.Out = os.Stdout
		return logger
	}

	logFile, err := os.OpenFile("plusy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = logFile
		multiWriters := io.MultiWriter(os.Stdout, logFile)
		logger.SetOutput(multiWriters)
	} else {
		logger.Warn("Failed to open logfile, defaulting to stderr")
	}

	return logger
}
