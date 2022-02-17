package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

const LOGFILE = "gogdl-ng.log"

var logger *logrus.Logger = nil

func NewLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.Formatter = &logrus.JSONFormatter{}
		logger.SetOutput(os.Stdout)

		file, err := os.OpenFile(LOGFILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

		if err != nil {
			logger.Fatal(err)
		}

		logger.SetOutput(file)
		logrus.RegisterExitHandler(func() {
			file.Close()
		})
	}

	return logger
}
