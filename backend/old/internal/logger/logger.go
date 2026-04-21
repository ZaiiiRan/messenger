package logger

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	instance *logrus.Logger
}

var loggerInstance *Logger
var once sync.Once

// Get instance of logger
func GetInstance() *Logger {
	once.Do(func() {
		logrusLogger := logrus.New()
		logrusLogger.SetFormatter(&logrus.JSONFormatter{})
		logrusLogger.SetOutput(os.Stdout)
		logrusLogger.SetLevel(logrus.DebugLevel)

		file, err := os.OpenFile("./app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			logrusLogger.SetOutput(file)
		} else {
			logrusLogger.Warn("Failed to log to file, defaulting to stdout")
		}

		loggerInstance = &Logger{
			instance: logrusLogger,
		}
	})
	return loggerInstance
}

// Log creation
func (l *Logger) Log(level logrus.Level, message string, action string, data interface{}, err error) {
	fields := logrus.Fields{
		"action": action,
		"data":   data,
	}

	if err != nil {
		fields["stackTrace"] = errors.WithStack(err).Error()
	}

	entry := l.instance.WithFields(fields)

	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	default:
		entry.Info(message)
	}
}

// Log error
func (l *Logger) Error(message string, action string, data interface{}, err error) {
	l.Log(logrus.ErrorLevel, message, action, data, err)
}

// Log fatal
func (l *Logger) Fatal(message string, action string) {
	l.Log(logrus.FatalLevel, message, action, nil, nil)
}
