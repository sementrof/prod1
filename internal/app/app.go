package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sementrof/prod1/internal/api"
	v1 "github.com/sementrof/prod1/internal/api/v1"
	Cassandra "github.com/sementrof/prod1/internal/cassandra"
	"github.com/sementrof/prod1/internal/kafka"
	"github.com/sementrof/prod1/internal/repository"
	"github.com/sementrof/prod1/internal/repository/postgres"
	"github.com/sirupsen/logrus"
)

func NewLogger(level logrus.Level, formatter logrus.Formatter) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	return logger
}
func Run() {

	logger := NewLogger(logrus.DebugLevel, &logrus.TextFormatter{FullTimestamp: true})
	if err := godotenv.Overload("env/.env"); err != nil {
		logger.Errorf("error loading env variables: %s", err)
		return
	}
	producer, err := kafka.InitKafka()
	if err != nil {
		logger.Errorf("Failed to initialize Kafka: %v", err)
		return
	}
	go kafka.StartLogConsumer()

	confPostgres := postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	logger.Debugf("config: %v", confPostgres)

	port := os.Getenv("PORT")

	dbConn, err := postgres.Connection(confPostgres)
	if err != nil {
		logger.Errorf("error connecting to database: %v", err)
		return
	}
	logger.Infof("Connected to PostgreSQL")
	fmt.Println(dbConn)

	Cass, erro := Cassandra.InitCassandra()
	if erro != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
		return
	}
	err = Cass.CreateLogTable()
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
		return
	}
	defer Cass.Close()

	repo := repository.NewRepository(dbConn, logger)
	apiserver := v1.NewServer(logger, repo, &producer, &Cass)
	router := api.SetupRouter(apiserver)

	go func() {
		// use https
		//  if err := http.ListenAndServeTLS(":"+port, "env/server.crt", "env/server.key", router); err != nil
		if err := http.ListenAndServe(":"+port, router); err != nil {
			logger.Errorf("Cann't run server: %v", err)
			return
		}
		logger.Infof("Server is running on port %s", port)
	}()
	logger.Infof("Rest server is running on port: %s", port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
