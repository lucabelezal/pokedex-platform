package httpclient

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/sony/gobreaker/v2"
)

// CircuitBreakerConfig define os parâmetros de resiliência.
type CircuitBreakerConfig struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	FailureCount  uint32
	RetryMax      int
	RetryBackoff  []time.Duration
}

// DefaultCircuitBreakerConfig retorna a configuração padrão.
func DefaultCircuitBreakerConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:         name,
		MaxRequests:  1,
		Interval:     60 * time.Second,
		Timeout:      30 * time.Second,
		FailureCount: 5,
		RetryMax:     3,
		RetryBackoff: []time.Duration{1 * time.Second, 3 * time.Second, 10 * time.Second},
	}
}

// CircuitBreakerClient é um wrapper que adiciona circuit breaker e retries ao http.Client.
type CircuitBreakerClient struct {
	name   string
	client *http.Client
	cb     *gobreaker.CircuitBreaker[*http.Response]
	retry  retryConfig
}

type retryConfig struct {
	max     int
	backoff []time.Duration
}

// NewCircuitBreakerClient cria um cliente HTTP com circuit breaker e retries.
func NewCircuitBreakerClient(inner *http.Client, cfg CircuitBreakerConfig) *CircuitBreakerClient {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.FailureCount
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			slog.Warn("circuit breaker mudou de estado",
				"name", name,
				"from", from.String(),
				"to", to.String(),
			)
		},
	}

	cb := gobreaker.NewCircuitBreaker[*http.Response](settings)

	return &CircuitBreakerClient{
		name:   cfg.Name,
		client: inner,
		cb:     cb,
		retry: retryConfig{
			max:     cfg.RetryMax,
			backoff: cfg.RetryBackoff,
		},
	}
}

// Do executa a requisição HTTP com circuit breaker e retries.
func (c *CircuitBreakerClient) Do(req *http.Request) (*http.Response, error) {
	return c.cb.Execute(func() (*http.Response, error) {
		return c.doWithRetry(req)
	})
}

func (c *CircuitBreakerClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retry.max; attempt++ {
		resp, err := c.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			resp.Body.Close()
			lastErr = &httpStatusError{StatusCode: resp.StatusCode}
		}

		if attempt >= c.retry.max {
			break
		}

		delay := c.backoffDuration(attempt)
		slog.Warn("retrying request",
			"name", c.name,
			"attempt", attempt+1,
			"delay", delay.String(),
			"error", lastErr.Error(),
		)
		time.Sleep(delay)
	}

	return nil, lastErr
}

func (c *CircuitBreakerClient) backoffDuration(attempt int) time.Duration {
	if attempt < len(c.retry.backoff) {
		base := c.retry.backoff[attempt]
		jitter := time.Duration(float64(base) * 0.1 * (rand.Float64()*2 - 1))
		return base + jitter
	}
	return 10 * time.Second
}

type httpStatusError struct {
	StatusCode int
}

func (e *httpStatusError) Error() string {
	return "http status " + http.StatusText(e.StatusCode)
}
