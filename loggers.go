package fairway

import (
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

func loggers() ([]byte, error) {
	if load {
		loggerData = LoggerData{Levels: levels}

		loggerData.Loggers = map[string]LogStruct{}
		loggerData.Loggers["ROOT"] = LogStruct{
			ConfiguredLevel: strings.ToUpper(LogLevel.String()),
			EffectiveLevel:  strings.ToUpper(LogLevel.String()),
		}
		load = false
	}
	if LogLevel.String() == "debug" {
		logger.Info("INFO")
		logger.Warn("WARN")
		logger.Error("ERROR")
		logger.Debug("DEBUG")
		logger.Trace("TRACE")
	}
	return toJson(loggerData), nil
}

func loggersUpdate(key string, update LogStruct) ([]byte, error) {
	setLog(update.ConfiguredLevel)
	update.EffectiveLevel = update.ConfiguredLevel
	loggerData.Loggers[key] = update
	return toJson(loggerData), nil
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
