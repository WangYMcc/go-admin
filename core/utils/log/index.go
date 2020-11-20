package log

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info *log.Logger
	Error *log.Logger
	WARN * log.Logger
)

func init(){
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	WARN = log.New(os.Stderr, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
}