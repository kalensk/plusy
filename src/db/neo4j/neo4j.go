package neo4j

import (
	"time"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/sirupsen/logrus"
)

type Neo4j struct {
	log               *logrus.Logger
	pool              golangNeo4jBoltDriver.DriverPool
	timeout           time.Duration
	lastOffsetPrinted int64
}

func New(log *logrus.Logger, databaseAddrPort string) (*Neo4j, error) {
	log.Info("Attempting to connect to Neo4j on ", databaseAddrPort)

	connString := "bolt://username:password@localhost:7687"
	maxConnections := 20
	pool, err := golangNeo4jBoltDriver.NewDriverPool(connString, maxConnections)

	timeout, err := time.ParseDuration("5s")
	if err != nil {
		return nil, err
	}

	return &Neo4j{log: log, pool: pool, timeout: timeout}, nil
}

func (n *Neo4j) Ping() error {
	conn, err := n.pool.OpenPool()
	if err != nil {
		return err
	}

	conn.PrepareNeo("")
	return nil
}

func (n *Neo4j) Close() error {
	// pool takes care of close...although look into using
	// golangNeo4jBoltDriver.NewClosableDriverPool()
	return nil
}

func (n *Neo4j) ClearData() (string, error) {
	// instead just create a new test database each time
	return "", nil
}
