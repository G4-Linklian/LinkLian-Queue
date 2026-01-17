package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func LogInfo(v ...interface{}) {
	InfoLogger.Println(v...)
}

func LogError(v ...interface{}) {
	ErrorLogger.Println(v...)
}

func LogInfof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

func LogErrorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}
