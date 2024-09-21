package wlog

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func createTextLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}
	return logger
}

func createStdoutLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
	return logger
}

func createStderrLogger() *logrus.Logger {
	logger := createStdoutLogger()
	logger.SetOutput(os.Stderr)
	return logger
}

func createDiscardLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   false,
		DisableColors: true,
		FullTimestamp: false,
	}
	logger.SetOutput(ioutil.Discard)
	return logger
}
