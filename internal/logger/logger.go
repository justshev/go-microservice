package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	z zerolog.Logger
}

func New(serviceName string,level string) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	z := zerolog.New(os.Stdout).
	Level(lvl).With().Str("service", serviceName).Timestamp().Logger()

	return &Logger{z: z}
	
}

func (l *Logger) Info(msg string)  { 
	l.z.Info().Msg(msg)
 }
func (l *Logger) Error(msg string) {
	l.z.Error().Msg(msg)
}

func (l *Logger) Raw() zerolog.Logger {
	return l.z
}