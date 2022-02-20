package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

const LOG_FILE = "gogdl-ng.log"

var logger *logrus.Logger = nil

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

func NewLogger1(logFileName string) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		return nil, err
	}

	logger.SetOutput(file)
	logrus.RegisterExitHandler(func() {
		file.Close()
	})

	return logger, nil
}

func NewLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		logger.Fatal(err)
	}

	logger.SetOutput(file)
	logrus.RegisterExitHandler(func() {
		file.Close()
	})

	return logger
}
