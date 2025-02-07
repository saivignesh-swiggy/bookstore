package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func Init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}

func Info(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(message)
}

func Error(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Error(message)
}

func Debug(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Debug(message)
}

func Fatal(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Fatal(message)
}
