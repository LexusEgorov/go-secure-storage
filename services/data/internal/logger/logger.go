package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init(env string) *logrus.Logger {
	log := logrus.New()

	log.Out = os.Stdout

	switch env {
	case "local":
		log.SetLevel(logrus.DebugLevel)
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	default:
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	return log
}
