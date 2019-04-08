package core

import (
	"io"
	"log"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogger(o io.Writer) {
	Info = log.New(o, "[Info]    ", log.Ldate|log.Ltime)
	Warning = log.New(o, "[Warning] ", log.Ldate|log.Ltime)
	Error = log.New(o, "[Error]   ", log.Ldate|log.Ltime)
}
