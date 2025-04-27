package kafka

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type LogMessage struct {
	Event   string `json:"event"`
	Message string `json:"message"`
	UserId  string `json:"user_id,omitempty"`
}

func StartLogConsumer() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9092,127.0.0.1:9093",
		"group.id":          "log-consumer-group10",
		"auto.offset.reset": "earliest",
		// "debug":              "cgrp,broker,fetch,protocol",
		"enable.auto.commit": "true",
	}
	logrus.Infof("Создание нового Kafka consumer с конфигурацией: %+v", config)
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		logrus.Errorf("Failed to create Kafka consumer: %v", err)
		return
	}
	defer func() {
		consumer.Close()
	}()

	// Колбэк для обработки ребалансировки
	rebalanceCb := func(c *kafka.Consumer, e kafka.Event) error {
		switch ev := e.(type) {
		case kafka.AssignedPartitions:
			if err := c.Assign(ev.Partitions); err != nil {
				logrus.Errorf("Ошибка назначения партиций: %v", err)
			}
		case kafka.RevokedPartitions:
			logrus.Infof("Партиции отозваны: %+v", ev.Partitions)
			if err := c.Unassign(); err != nil {
				logrus.Errorf("Ошибка отзыва партиций: %v", err)
			}
		default:
			logrus.Infof("Ребалансировочное событие: %v", e)
		}
		return nil
	}

	err = consumer.SubscribeTopics([]string{"logs"}, rebalanceCb)
	if err != nil {
		logrus.Errorf("Failed to subscribe to Kafka topics: %v", err)
		return
	}

	logFile, err := os.OpenFile("monitoring/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("Failed to open log file: %v", err)
		return
	}
	defer func() {
		logFile.Close()
	}()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// Основной цикл чтения сообщений с таймаутом, чтобы не блокироваться бесконечно
	for {
		msg, err := consumer.ReadMessage(1000)
		if err != nil {
			kErr, ok := err.(kafka.Error)
			if ok && kErr.Code() == kafka.ErrTimedOut {
				logrus.Debugf("Таймаут ожидания сообщения")
				continue
			}
			logrus.Errorf("Ошибка чтения сообщения из Kafka: %v", err)
			continue
		}
		logrus.Infof("Получено сообщение: %s", string(msg.Value))

		var logEntry LogMessage
		if err := json.Unmarshal(msg.Value, &logEntry); err != nil {
			continue
		}

		logLine := logEntry.Event + ": " + logEntry.Message
		if logEntry.UserId != "" {
			logLine += " (User ID: " + logEntry.UserId + ")"
		}
		logLine += "\n"

		if _, err := logFile.WriteString(logLine); err != nil {
			logrus.Errorf("Ошибка записи в файл: %v", err)
		} else {
			logrus.Infof("Лог записан в файл: %s", logLine)
		}
	}
}
