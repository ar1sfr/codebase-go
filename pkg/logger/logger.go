package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(env string) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return fmt.Sprintf("| %-6s|", i)
		},
	}

	level := zerolog.InfoLevel
	if env == "development" {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	log = zerolog.New(output).With().Timestamp().Caller().Logger()
}

func Debug(message string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Debug().Fields(fieldsToMap(fields...)).Msg(message)
		return
	}
	log.Debug().Msg(message)
}

func Info(message string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Info().Fields(fieldsToMap(fields...)).Msg(message)
		return
	}
	log.Info().Msg(message)
}

func Warn(message string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Warn().Fields(fieldsToMap(fields...)).Msg(message)
		return
	}
	log.Warn().Msg(message)
}

func Error(message string, err error, fields ...interface{}) {
	if len(fields) > 0 {
		log.Error().Err(err).Fields(fieldsToMap(fields...)).Msg(message)
		return
	}
	log.Error().Err(err).Msg(message)
}

func Fatal(message string, err error, fields ...interface{}) {
	if len(fields) > 0 {
		log.Fatal().Err(err).Fields(fieldsToMap(fields...)).Msg(message)
		return
	}
	log.Fatal().Err(err).Msg(message)
}

func fieldsToMap(fields ...interface{}) map[string]interface{} {
	if len(fields)%2 != 0 {
		return nil
	}

	fieldMap := make(map[string]interface{}, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		fieldMap[key] = fields[i+1]
	}

	return fieldMap
}
