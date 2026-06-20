package zap_betterstack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Custom log level for NOTICE (below DebugLevel, non-error informational logs)
const NoticeLevel zapcore.Level = -2

//logEntry represents a single log entry for Better Stack
type logEntry struct {
	Timestamp  string         `json:"timestamp"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	TraceID    string         `json:"traceID"`
	Attributes map[string]any `json:"attributes"`
}

//BetterStackLoggingMiddleware captures HTTP request metadata and logs it based on environment
func BetterStackLoggingMiddleware(sourceToken, environment, BetterStackUploadURL string, logger *zap.Logger) gin.HandlerFunc {
	client := &http.Client{Timeout: 10 * time.Second}

	//File writer for development environment
	var fileWriter io.Writer
	var fileMu sync.Mutex
	if environment == "development" {
		f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Error("Failed to open log file", zap.Error(err))
			fileWriter = os.Stderr
		} else {
			fileWriter = f
		}
		//Combine file and stdout for development
		fileWriter = io.MultiWriter(fileWriter, os.Stdout)
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		//Generate TraceID and set in header and context
		traceID := uuid.New().String()
		c.Request.Header.Set("X-Trace-ID", traceID)
		//c.Set("TraceID", traceID)

		//Process request
		c.Next()

		//Determine log level based on response status
		status := c.Writer.Status()
		var levelStr string

		switch {
		case status >= 500:
			//server errors
			levelStr = "ERROR"
		case status >= 400:
			//client errors
			levelStr = "WARN"
		case status >= 300:
			//redirects
			levelStr = "INFO"
		case status >= 200:
			//success
			levelStr = "NOTICE"
		default:
			//unknown cases
			levelStr = "DEBUG"
		}

		if c.Request.Method == "OPTIONS" {
			return
		}

		//Create log entry
		latency := time.Since(start)
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		entry := logEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Level:     levelStr,
			Message:   "HTTP request",
			TraceID:   traceID,
			Attributes: map[string]any{
				"method":     c.Request.Method,
				"path":       path,
				"query":      query,
				"status":     status,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
				"latency_ms": latency.Milliseconds(),
				"body":       string(bodyBytes),
			},
		}

		//Marshal log to JSON
		body, err := json.Marshal(entry)
		if err != nil {
			logger.Error("Failed to marshal log", zap.Error(err))
			return
		}

		if environment == "development" {
			//Write to file only in development (avoid console writes)
			fileMu.Lock()
			_, err := fileWriter.Write(append(body, '\n'))
			fileMu.Unlock()
			if err != nil {
				logger.Error("Failed to write log to file", zap.Error(err))
			}
			return
		}

		//Production: Send log to Better Stack
		req, err := http.NewRequest("POST", BetterStackUploadURL, bytes.NewReader(body))
		if err != nil {
			//logger.Error("Failed to create HTTP request", zap.Error(err))
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+sourceToken)

		//Send log asynchronously
		go func() {
			resp, err := client.Do(req)
			if err != nil {
				logger.Error("Failed to send log to Better Stack", zap.Error(err))
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusAccepted {
				logger.Error("Unexpected response from Better Stack", zap.String("status", resp.Status))
			}
		}()
	}
}
