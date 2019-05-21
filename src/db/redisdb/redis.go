package redisdb

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	log               *logrus.Logger
	conn              redis.Conn
	timeout           time.Duration
	lastOffsetPrinted int64
}

func New(log *logrus.Logger, databaseAddrPort string) (*Redis, error) {
	log.Info("Attempting to connect to Redis on ", databaseAddrPort)
	conn, err := redis.Dial("tcp", databaseAddrPort)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration("5s")
	if err != nil {
		return nil, err
	}

	return &Redis{log: log, conn: conn, timeout: timeout}, nil
}

func (r *Redis) Ping() error {
	_, err := redis.String(r.conn.Do("Ping"))
	return err
}

func (r *Redis) Close() error {
	return r.conn.Close()
}

func (r *Redis) ClearData() (string, error) {
	return redis.String(r.conn.Do("FLUSHALL"))
}
