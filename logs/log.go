package logs

// handles all the logs of the applications

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger interface
type LoggerInterface interface {
	Info(message string)
	Warning(message string)
	Critical(message string)
	Debug(message string)
	Route(message string)
}

// Logger struct
type Logger struct {
}

// NewLogger creates a new Logger instance and ensures the logs directory exists.
func NewLogger() LoggerInterface {
	// Create logs directory if not exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &Logger{}
}

// createLogFile opens the log file in append mode and returns the file.
func (l *Logger) createLogFile(filename string) *os.File {
	// Open log file with append mode
	file, err := os.OpenFile("logs/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// writeLog writes a log message with the specified level, color, and file name.
func (l *Logger) writeLog(level, color, message, filename string) {
	// Get the time to log
	timeStamp := time.Now().Format("2006-01-02 15:04:05")

	// Create the log file for the level
	file := l.createLogFile(filename)
	defer file.Close()

	// Set color for output
	var colorCode string
	switch color {
	case "green":
		colorCode = "\033[32m"
	case "red":
		colorCode = "\033[31m"
	case "yellow":
		colorCode = "\033[33m"
	case "blue":
		colorCode = "\033[34m"
	default:
		colorCode = "\033[0m"
	}

	// Format the log message
	logMessage := fmt.Sprintf("%s%s[%s] %s\033[0m\n", colorCode, timeStamp, level, message)

	// Output to console
	fmt.Print(logMessage)

	// Write to log file
	log.SetOutput(file)
	log.Println(logMessage)
}

// Info logs an info message to both the terminal and info.log.
func (l *Logger) Route(message string) {
	l.writeLog("INFO", "blue", message, "route.log")
}

func (l *Logger) Info(message string) {
	l.writeLog("INFO", "green", message, "info.log")
}

// Warning logs a warning message to both the terminal and warning.log.
func (l *Logger) Warning(message string) {
	l.writeLog("WARNING", "yellow", message, "warning.log")
}

// Critical logs a critical message to both the terminal and critical.log.
func (l *Logger) Critical(message string) {
	l.writeLog("CRITICAL", "red", message, "critical.log")
}

// Debug logs a debug message to both the terminal and debug.log.
func (l *Logger) Debug(message string) {
	l.writeLog("DEBUG", "blue", message, "debug.log")
}

// Main function to demonstrate usage of the logger
func main() {
	// Initialize Logger
	var logger LoggerInterface = NewLogger()

	// Example Logs
	logger.Info("This is an info message")
	logger.Warning("This is a warning message")
	logger.Critical("This is a critical message")
	logger.Debug("This is a debug message for easier debugging")
}
