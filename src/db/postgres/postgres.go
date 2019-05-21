package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	log               *logrus.Logger
	db                *sql.DB
	timeout           time.Duration
	lastOffsetPrinted int64
}

func New(log *logrus.Logger, databaseAddrPort string) (*Postgres, error) {
	log.Info("Attempting to connect to Postgres on ", databaseAddrPort)

	pq.QuoteIdentifier("turtles") // just to get lib/pq import

	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration("5s")
	if err != nil {
		return nil, err
	}

	return &Postgres{log: log, db: db, timeout: timeout}, nil
}

func (p *Postgres) Ping() error {
	return p.db.Ping()
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) ClearData() (string, error) {
	// instead just create a new test database each time
	return "", nil
}
