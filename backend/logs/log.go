package logs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	INFO  = "INFO"
	DEBUG = "DEBUG"
	ERROR = "ERROR"
	HTTP  = "HTTP"
)

type Logger struct {
	file *os.File
}


// NewLogger initializes a logger with a daily log file
func NewLogger() (*Logger, error) {
	// Ensure the logs directory exists
	logDir := "logs"
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate the log filename based on the current date
	logFilename := filepath.Join(logDir, time.Now().Format("02-01-06.log"))
	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{file: logFile}, nil
}

// logMessage formats and writes the log message to both file and console
func (l *Logger) logMessage(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("%s - %s - %s\n", timestamp, level, message)

	// Write to the log file
	_, _ = l.file.WriteString(logMessage)

	// Print to the console
	fmt.Print(logMessage)
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.logMessage(INFO, message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.logMessage(DEBUG, message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.logMessage(ERROR, message)
}

// HTTP logs an HTTP request message
func (l *Logger) HTTP(message string, ip string) {
	l.logMessage(HTTP, message)
}

// Close releases resources held by the logger
func (l *Logger) Close() {
	_ = l.file.Close()
}

// Middleware to log HTTP requests with IP
func LoggingMiddleware(logger *Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from the request
		ip := getIPAddress(r)
		logger.HTTP(fmt.Sprintf("HTTP request: %s %s", r.Method, r.URL.Path), ip)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Extract the IP address from the HTTP request
func getIPAddress(r *http.Request) string {
	// Check for the X-Forwarded-For header (used by proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Return the first IP in the list (client IP)
		return strings.Split(xff, ",")[0]
	}

	// Fall back to the remote address
	ip := r.RemoteAddr
	// If the address contains a port, remove it
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	if ip == "::1" {
		return "127.0.0.1" // Map IPv6 loopback to IPv4 loopback
	}

	return ip
}
