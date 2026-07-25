package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpclient "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircuitBreakerClientSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := httpclient.DefaultCircuitBreakerConfig("test-success")
	cfg.FailureCount = 2
	inner := &http.Client{Timeout: 2 * time.Second}
	cb := httpclient.NewCircuitBreakerClient(inner, cfg)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	resp, err := cb.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestCircuitBreakerOpensAfterFailures(t *testing.T) {
	failureCount := uint32(2)
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := httpclient.DefaultCircuitBreakerConfig("test-failures")
	cfg.FailureCount = failureCount
	cfg.MaxRequests = 0
	inner := &http.Client{Timeout: 2 * time.Second}
	cb := httpclient.NewCircuitBreakerClient(inner, cfg)

	for i := 0; i < int(failureCount+2); i++ {
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		resp, err := cb.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
		if i >= int(failureCount+1) { // após abrir + 1 tentativa extra
			assert.Error(t, err)
		}
	}

	assert.GreaterOrEqual(t, callCount, int(failureCount), "deve ter feito pelo menos N chamadas antes de abrir")
}

func TestCircuitBreakerRetryOn5xx(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := httpclient.DefaultCircuitBreakerConfig("test-retry")
	cfg.RetryMax = 3
	cfg.RetryBackoff = []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond}
	cfg.FailureCount = 10 // não deve abrir com tão poucas falhas
	inner := &http.Client{Timeout: 2 * time.Second}
	cb := httpclient.NewCircuitBreakerClient(inner, cfg)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, err := cb.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	assert.GreaterOrEqual(t, attempts, 3, "deve ter retentado pelo menos 3 vezes")
}

func TestCircuitBreakerNoRetryOn4xx(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := httpclient.DefaultCircuitBreakerConfig("test-noretry-4xx")
	cfg.RetryMax = 3
	cfg.RetryBackoff = []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond}
	inner := &http.Client{Timeout: 2 * time.Second}
	cb := httpclient.NewCircuitBreakerClient(inner, cfg)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, err := cb.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
	assert.Equal(t, 1, attempts, "não deve retentar em erros 4xx")
}

func TestCircuitBreakerTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer server.Close()

	cfg := httpclient.DefaultCircuitBreakerConfig("test-timeout")
	cfg.FailureCount = 1
	cfg.RetryMax = 0
	inner := &http.Client{Timeout: 10 * time.Millisecond}
	cb := httpclient.NewCircuitBreakerClient(inner, cfg)

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := cb.Do(req)
	assert.Error(t, err, "deve falhar por timeout")
}

func TestCircuitBreakerConfiguration(t *testing.T) {
	cfg := httpclient.DefaultCircuitBreakerConfig("test-config")

	assert.Equal(t, "test-config", cfg.Name)
	assert.Equal(t, uint32(1), cfg.MaxRequests)
	assert.Equal(t, 60*time.Second, cfg.Interval)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
	assert.Equal(t, uint32(5), cfg.FailureCount)
	assert.Equal(t, 3, cfg.RetryMax)
	assert.Len(t, cfg.RetryBackoff, 3)
}

func TestHTTPStatusError(t *testing.T) {
	t.Run("5xx error triggers retry and message", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := httpclient.DefaultCircuitBreakerConfig("test-status-error")
		cfg.RetryMax = 0
		cfg.FailureCount = 10
		inner := &http.Client{Timeout: 2 * time.Second}
		cb := httpclient.NewCircuitBreakerClient(inner, cfg)

		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		_, err := cb.Do(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Internal Server Error")
	})
}
