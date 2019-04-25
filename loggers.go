package fairway

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
)

var levels = []string{
	"OFF",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

var loggerData = LoggerData{}
var load = true

type LoggerData struct {
	Levels  []string             `json:"levels"`
	Loggers map[string]LogStruct `json:"loggers"`
}

type LogStruct struct {
	ConfiguredLevel string `json:"configuredLevel"`
	EffectiveLevel  string `json:"effectiveLevel"`
}

func loggers(key string) ([]byte, error) {
	if load {
		loggerData = LoggerData{Levels: levels}

		loggerData.Loggers = map[string]LogStruct{}
		loggerData.Loggers["ROOT"] = LogStruct{
			ConfiguredLevel: strings.ToUpper(LogLevel.String()),
			EffectiveLevel:  strings.ToUpper(LogLevel.String()),
		}
		load = false
	}
	if key != "" {
		return toJson(loggerData.Loggers[key]), nil
	}

	return toJson(loggerData), nil
}

func loggersUpdate(key string, update LogStruct) ([]byte, error) {

	setLog(update.ConfiguredLevel)
	update.EffectiveLevel = update.ConfiguredLevel
	if _, ok := loggerData.Loggers[key]; ok {
		loggerData.Loggers[key] = update
		return toJson(loggerData.Loggers[key]), nil
	}
	return toJson(LoggerData{}), errors.New("Key not found")
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
