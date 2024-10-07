package logger

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Logger wraps the standard log.Logger and adds Slack integration
type Logger struct {
	*log.Logger
	slackWebhookURL string
	slackLevels     map[string]bool // Tracks which log levels should send to Slack
	mu              sync.Mutex
}

// NewLogger creates a new logger with Slack integration
func NewLogger(slackWebhookURL string, slackLevels []string) *Logger {
	slackLevelsMap := make(map[string]bool)
	for _, level := range slackLevels {
		slackLevelsMap[level] = true
	}

	return &Logger{
		Logger:          log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		slackWebhookURL: slackWebhookURL,
		slackLevels:     slackLevelsMap,
	}
}

// Println logs a message and sends to Slack if the level is enabled for Slack
func (l *Logger) Println(level string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Log to standard output
	l.Logger.Println(level+":", v)

	// Send to Slack if the level is enabled
	if l.slackLevels[level] {
		go l.sendToSlack(level, v...)
	}
}

// Printf logs a formatted message and sends to Slack if the level is enabled for Slack
func (l *Logger) Printf(level string, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Log to standard output
	l.Logger.Printf(level+": "+format, v...)

	// Send to Slack if the level is enabled
	if l.slackLevels[level] {
		go l.sendToSlack(level, v...)
	}
}

// sendToSlack sends the log message to the configured Slack webhook
func (l *Logger) sendToSlack(level string, v ...interface{}) {
	slackPayload := bytes.NewBuffer([]byte(
		`{"text": "` + level + `: ` + formatMessage(v...) + `"}`))

	req, err := http.NewRequest("POST", l.slackWebhookURL, slackPayload)
	if err != nil {
		l.Logger.Println("Failed to create Slack request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		l.Logger.Println("Failed to send log to Slack:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Logger.Printf("Slack responded with status code %d", resp.StatusCode)
	}
}

// formatMessage formats the log message to be Slack-compatible
func formatMessage(v ...interface{}) string {
	return fmt.Sprint(v...)
}
