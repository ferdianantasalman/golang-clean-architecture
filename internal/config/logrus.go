package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	level, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logrus.Level(level))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
