package Cassandra

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type Cassandra interface {
	CreateLogTable() error
	InsertLog(id gocql.UUID, status string, text string) error
	GetLogByID(id gocql.UUID) (string, int, error)
	Close()
}

type CassandraWrapper struct {
	session *gocql.Session
}

func CreateKeyspace() error {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		logrus.Errorf("Ошибка подключения к Cassandra при создании keyspace: %s", err)
		return err
	}
	defer session.Close()

	query := `
	CREATE KEYSPACE IF NOT EXISTS myapp 
	WITH replication = {
		'class': 'SimpleStrategy',
		'replication_factor': 1
	}`
	return session.Query(query).Exec()
}

// func InitCassandra() (Cassandra, error) {
// 	cluster := gocql.NewCluster("127.0.0.1")
// 	cluster.Keyspace = "system"
// 	cluster.Consistency = gocql.Quorum

// 	sessions, err := cluster.CreateSession()
// 	if err != nil {
// 		logrus.Errorf("Не удалось подключиться к Cassandra: %s", err)
// 		return nil, err
// 	}
// 	return &CassandraWrapper{session: sessions}, nil
// }

func InitCassandra() (Cassandra, error) {
	const keyspace = "myapp"

	// Шаг 1: Подключаемся без keyspace
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum

	tempSession, err := cluster.CreateSession()
	if err != nil {
		logrus.Errorf("Ошибка подключения к Cassandra для создания keyspace: %s", err)
		return nil, err
	}
	defer tempSession.Close()

	// Шаг 2: Создаём keyspace, если он ещё не существует
	query := `
	CREATE KEYSPACE IF NOT EXISTS ` + keyspace + ` 
	WITH replication = {
		'class': 'SimpleStrategy',
		'replication_factor': 1
	}`
	if err := tempSession.Query(query).Exec(); err != nil {
		logrus.Errorf("Ошибка создания keyspace: %s", err)
		return nil, err
	}

	// Шаг 3: Подключаемся уже с нужным keyspace
	cluster.Keyspace = keyspace
	mainSession, err := cluster.CreateSession()
	if err != nil {
		logrus.Errorf("Ошибка подключения к Cassandra с keyspace '%s': %s", keyspace, err)
		return nil, err
	}

	return &CassandraWrapper{session: mainSession}, nil
}

func (c *CassandraWrapper) CreateLogTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		status TEXT,
		text TEXT
	)`
	return c.session.Query(query).Exec()
}

func (c *CassandraWrapper) InsertLog(id gocql.UUID, status string, text string) error {
	query := `INSERT INTO users (id, status, text) VALUES (?, ?, ?)`
	return c.session.Query(query, id, status, text).Exec()
}

func (c *CassandraWrapper) GetLogByID(id gocql.UUID) (string, int, error) {
	var status string
	var text int
	query := `SELECT status, text FROM users WHERE id = ? LIMIT 1`
	if err := c.session.Query(query, id).Scan(&status, &text); err != nil {
		return "", 0, err
	}
	return status, text, nil
}

func (c *CassandraWrapper) Close() {
	c.session.Close()
}
