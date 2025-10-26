package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LokiLogger struct {
	lokiURL string
	client  *http.Client
}

type LokiRequest struct {
	Streams []LokiStream `json:"streams"`
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func NewLokiLogger(lokiURL string) *LokiLogger {
	return &LokiLogger{
		lokiURL: lokiURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (l *LokiLogger) Info(ctx context.Context, msg string, fields ...shared.Field) {
	l.sendToLoki(ctx, "INFO", msg, fields...)
}

func (l *LokiLogger) Error(ctx context.Context, msg string, fields ...shared.Field) {
	l.sendToLoki(ctx, "ERROR", msg, fields...)
}

func (l *LokiLogger) Debug(ctx context.Context, msg string, fields ...shared.Field) {
	l.sendToLoki(ctx, "DEBUG", msg, fields...)
}

func (l *LokiLogger) sendToLoki(ctx context.Context, level, msg string, fields ...shared.Field) {
	if l.lokiURL == "" {
		return
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	logLine := fmt.Sprintf("%s: %s", level, msg)

	var requestID string
	for _, field := range fields {
		logLine += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		if field.Key == "request_id" {
			requestID, _ = field.Value.(string)
		}
	}

	if requestID == "" {
		const requestIDKey = "request_id"
		requestID, _ = ctx.Value(requestIDKey).(string)
	}

	if requestID == "" {
		requestID = "none"
	}

	stream := LokiStream{
		Stream: map[string]string{
			"job":        "snipme-api",
			"level":      level,
			"request_id": requestID,
		},
		Values: [][]string{
			{timestamp, logLine},
		},
	}
	if requestID != "" {
		stream.Stream["request_id"] = requestID
	}

	req := LokiRequest{
		Streams: []LokiStream{stream},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return
	}

	go func() {
		resp, err := l.client.Post(l.lokiURL+"/loki/api/v1/push", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error enviando log a Loki: %v\n", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			fmt.Printf("Loki respondió con status: %v\n", resp.Status)
		}
	}()

}
