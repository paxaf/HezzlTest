package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var globalLogger *Logger

type LogLevel int

const (
	DebugLog LogLevel = iota
	InfoLog
	WarnLog
	ErrorLog
	FatalLog
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(getLogLevel(level))

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()

	l := &Logger{
		logger: &logger,
	}

	globalLogger = l
	return l
}

func initLoggerIfNeeded(level LogLevel) {
	if globalLogger != nil {
		return
	}

	defaultLogLevel := "info"
	switch level {
	case DebugLog:
		defaultLogLevel = "debug"
	case WarnLog:
		defaultLogLevel = "warn"
	case ErrorLog:
		defaultLogLevel = "error"
	case FatalLog:
		defaultLogLevel = "fatal"
	}

	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		defaultLogLevel = envLevel
	}

	New(defaultLogLevel)
}

func Debug(message interface{}, args ...interface{}) {
	initLoggerIfNeeded(DebugLog)
	globalLogger.Debug(message, args...)
}

func Info(message string, args ...interface{}) {
	initLoggerIfNeeded(InfoLog)
	globalLogger.Info(message, args...)
}

func Warn(message string, args ...interface{}) {
	initLoggerIfNeeded(WarnLog)
	globalLogger.Warn(message, args...)
}

func Error(message interface{}, args ...interface{}) {
	initLoggerIfNeeded(ErrorLog)
	globalLogger.Error(message, args...)
}

func Fatal(message interface{}, args ...interface{}) {
	initLoggerIfNeeded(FatalLog)
	globalLogger.Fatal(message, args...)
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.log(DebugLog, message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log(InfoLog, message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(WarnLog, message, args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.log(ErrorLog, message, args...)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.log(FatalLog, message, args...)
}

func (l *Logger) log(level LogLevel, message interface{}, args ...interface{}) {
	var msgStr string
	switch v := message.(type) {
	case error:
		msgStr = v.Error()
	case string:
		msgStr = v
	default:
		msgStr = fmt.Sprintf("%v", v)
	}

	event := l.createEvent(level)
	if event == nil {
		return
	}

	switch len(args) {
	case 0:
		event.Msg(msgStr)
	case 1:
		if err, ok := args[0].(error); ok {
			event.Err(err).Msg(msgStr)
			return
		}
		if fields, ok := args[0].(map[string]interface{}); ok {
			event.Fields(fields).Msg(msgStr)
			return
		}
		if strings.Contains(msgStr, "%") {
			event.Msgf(msgStr, args[0])
		} else {
			event.Interface("arg", args[0]).Msg(msgStr)
		}
	default:
		if strings.Contains(msgStr, "%") {
			event.Msgf(msgStr, args...)
		} else {

			fields := make(map[string]interface{})
			for i, arg := range args {
				fields[fmt.Sprintf("arg%d", i)] = arg
			}
			event.Fields(fields).Msg(msgStr)
		}
	}
}

func (l *Logger) createEvent(level LogLevel) *zerolog.Event {
	switch level {
	case DebugLog:
		return l.logger.Debug()
	case InfoLog:
		return l.logger.Info()
	case WarnLog:
		return l.logger.Warn()
	case ErrorLog:
		return l.logger.Error()
	case FatalLog:
		return l.logger.Fatal()
	default:
		return l.logger.Info()
	}
}

func getLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
