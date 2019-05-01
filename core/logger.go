package core

import (
	"io"
	"log"
)

// Logger ...
type Logger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// InitLogger ...
func InitLogger(o io.Writer) *Logger {
	return &Logger{
		Info:    log.New(o, "[Info]    ", log.Ldate|log.Ltime),
		Warning: log.New(o, "[Warning] ", log.Ldate|log.Ltime),
		Error:   log.New(o, "[Error]   ", log.Ldate|log.Ltime),
	}
}
