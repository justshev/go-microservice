package logger

import (
	"log"
	"os"
)

type Logger struct {
	prefix string
}

func New(serviceName string) *Logger {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	return &Logger{prefix: serviceName}
}

func (l *Logger) Info(msg string)  { log.Printf("[INFO] [%s] %s", l.prefix, msg) }
func (l *Logger) Error(msg string) { log.Printf("[ERROR] [%s] %s", l.prefix, msg) }
