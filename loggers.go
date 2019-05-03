package fairway

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

var levels = []string{
	"OFF",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

var loggerData = loggerStruct{}
var load = true

type loggerStruct struct {
	Levels  []string             `json:"levels"`
	Loggers map[string]logStruct `json:"loggers"`
}

type logStruct struct {
	ConfiguredLevel string `json:"configuredLevel"`
	EffectiveLevel  string `json:"effectiveLevel"`
}

func loggers(key string) ([]byte, error) {
	if load {
		loggerData = loggerStruct{Levels: levels}

		loggerData.Loggers = map[string]logStruct{}
		loggerData.Loggers["ROOT"] = logStruct{
			ConfiguredLevel: strings.ToUpper(logLevel.String()),
			EffectiveLevel:  strings.ToUpper(logLevel.String()),
		}
		load = false
	}
	if key != "" {
		return toJSON(loggerData.Loggers[key]), nil
	}

	return toJSON(loggerData), nil
}

func loggersUpdate(key string, update logStruct) ([]byte, error) {

	setLog(update.ConfiguredLevel)
	update.EffectiveLevel = update.ConfiguredLevel
	if _, ok := loggerData.Loggers[key]; ok {
		loggerData.Loggers[key] = update
		return toJSON(loggerData.Loggers[key]), nil
	}
	return toJSON(loggerStruct{}), errors.New("Key not found")
}

func setLog(level string) {
	switch level {
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		logger.SetLevel(logrus.WarnLevel)
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
	case "ERROR":
		logger.SetLevel(logrus.DebugLevel)
	case "TRACE":
		logger.SetLevel(logrus.DebugLevel)
	}
}
