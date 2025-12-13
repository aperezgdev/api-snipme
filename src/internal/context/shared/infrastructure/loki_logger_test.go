package infrastructure

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestLokiLogger_SendsLogsToLoki(t *testing.T) {
	var receivedRequests []LokiRequest
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/loki/api/v1/push", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var lokiReq LokiRequest
		err = json.Unmarshal(body, &lokiReq)
		assert.NoError(t, err)

		mu.Lock()
		receivedRequests = append(receivedRequests, lokiReq)
		mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	logger := NewLokiLogger(server.URL)
	ctx := context.Background()

	t.Run("Info log is sent correctly", func(t *testing.T) {
		logger.Info(ctx, "test info message", shared.NewField("key", "value"))
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		assert.Greater(t, len(receivedRequests), 0)
		lastReq := receivedRequests[len(receivedRequests)-1]
		assert.Equal(t, "INFO", lastReq.Streams[0].Stream["level"])
		assert.Contains(t, lastReq.Streams[0].Values[0][1], "test info message")
		assert.Contains(t, lastReq.Streams[0].Values[0][1], "key=value")
	})

	t.Run("Error log is sent correctly", func(t *testing.T) {
		logger.Error(ctx, "test error message", shared.NewField("error", "something went wrong"))
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		lastReq := receivedRequests[len(receivedRequests)-1]
		assert.Equal(t, "ERROR", lastReq.Streams[0].Stream["level"])
		assert.Contains(t, lastReq.Streams[0].Values[0][1], "test error message")
	})

	t.Run("Debug log is sent correctly", func(t *testing.T) {
		logger.Debug(ctx, "test debug message")
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		lastReq := receivedRequests[len(receivedRequests)-1]
		assert.Equal(t, "DEBUG", lastReq.Streams[0].Stream["level"])
		assert.Contains(t, lastReq.Streams[0].Values[0][1], "test debug message")
	})
}

func TestLokiLogger_WithoutURL_DoesNotSendLogs(t *testing.T) {
	t.Parallel()
	logger := NewLokiLogger("")
	ctx := context.Background()

	logger.Info(ctx, "test message")
	logger.Error(ctx, "test error")
	logger.Debug(ctx, "test debug")
}

func TestLokiLogger_WithRequestID(t *testing.T) {
	var receivedRequests []LokiRequest
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var lokiReq LokiRequest
		json.Unmarshal(body, &lokiReq)

		mu.Lock()
		receivedRequests = append(receivedRequests, lokiReq)
		mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	logger := NewLokiLogger(server.URL)
	ctx := context.WithValue(context.Background(), "request_id", "test-req-123")

	logger.Info(ctx, "test message with request id")
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Greater(t, len(receivedRequests), 0)
	lastReq := receivedRequests[len(receivedRequests)-1]
	assert.Equal(t, "test-req-123", lastReq.Streams[0].Stream["request_id"])
}

func TestLokiLogger_HandlesServerError(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	logger := NewLokiLogger(server.URL)
	ctx := context.Background()

	logger.Info(ctx, "test message")
	time.Sleep(100 * time.Millisecond)
}

func TestLokiLogger_HandlesInvalidURL(t *testing.T) {
	t.Parallel()
	logger := NewLokiLogger("http://invalid-url-that-does-not-exist:9999")
	ctx := context.Background()

	logger.Info(ctx, "test message")
	time.Sleep(100 * time.Millisecond)
}
