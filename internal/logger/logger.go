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

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()

	l := &Logger{
		logger: &logger,
	}

	globalLogger = l
	return l
}

func Debug(message interface{}, args ...interface{}) {
	if globalLogger == nil {
		defaultLogLevel := "debug"
		if os.Getenv("LOG_LEVEL") != "" {
			defaultLogLevel = os.Getenv("LOG_LEVEL")
		}
		New(defaultLogLevel)
	}
	globalLogger.Debug(message, args...)
}

func Info(message string, args ...interface{}) {
	if globalLogger == nil {
		defaultLogLevel := "info"
		if os.Getenv("LOG_LEVEL") != "" {
			defaultLogLevel = os.Getenv("LOG_LEVEL")
		}
		New(defaultLogLevel)
	}
	globalLogger.Info(message, args...)
}

func Warn(message string, args ...interface{}) {
	if globalLogger == nil {
		defaultLogLevel := "warn"
		if os.Getenv("LOG_LEVEL") != "" {
			defaultLogLevel = os.Getenv("LOG_LEVEL")
		}
		New(defaultLogLevel)
	}
	globalLogger.Warn(message, args...)
}

func Error(message interface{}, args ...interface{}) {
	if globalLogger == nil {
		defaultLogLevel := "error"
		if os.Getenv("LOG_LEVEL") != "" {
			defaultLogLevel = os.Getenv("LOG_LEVEL")
		}
		New(defaultLogLevel)
	}
	globalLogger.Error(message, args...)
}

func Fatal(message interface{}, args ...interface{}) {
	if globalLogger == nil {
		defaultLogLevel := "info"
		if os.Getenv("LOG_LEVEL") != "" {
			defaultLogLevel = os.Getenv("LOG_LEVEL")
		}
		New(defaultLogLevel)
	}
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
	var msg string
	switch v := message.(type) {
	case error:
		msg = v.Error()
	case string:
		msg = v
	default:
		msg = fmt.Sprintf("%v", v)
	}

	if len(args) == 1 {
		if fields, ok := args[0].(map[string]interface{}); ok {
			l.msgWithFields(level, msg, fields)
			return
		}
	}

	if len(args) > 0 && strings.Contains(msg, "%") {
		msg = fmt.Sprintf(msg, args...)
	}

	l.msg(level, msg)
}

func (l *Logger) msgWithFields(level LogLevel, message string, fields map[string]interface{}) {
	event := l.logger.With().Fields(fields).Logger()
	switch level {
	case DebugLog:
		event.Debug().Msg(message)
	case InfoLog:
		event.Info().Msg(message)
	case WarnLog:
		event.Warn().Msg(message)
	case ErrorLog:
		event.Error().Msg(message)
	case FatalLog:
		event.Fatal().Msg(message)
	}
}

func (l *Logger) msg(level LogLevel, message interface{}) {
	switch level {
	case DebugLog:
		l.logger.Debug().Msgf("%v", message)
	case InfoLog:
		l.logger.Info().Msgf("%v", message)
	case WarnLog:
		l.logger.Warn().Msgf("%v", message)
	case ErrorLog:
		l.logger.Error().Msgf("%v", message)
	case FatalLog:
		l.logger.Fatal().Msgf("%v", message)
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
