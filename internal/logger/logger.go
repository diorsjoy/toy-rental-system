package logger

import (
	"log"
	"os"
)

var (
	// Define the logger
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func Init() {
	// Create a new log file or append to the existing one
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}

	// Initialize the loggers
	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
