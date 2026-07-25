package httphandler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisPrefix = "ratelimit:"

type rateLimiter interface {
	Allow(ctx context.Context, clientID string) bool
}

type redisRateLimiter struct {
	client      *redis.Client
	maxRequests int
	window      time.Duration
}

func newRedisRateLimiter(client *redis.Client, maxRequests int, window time.Duration) *redisRateLimiter {
	return &redisRateLimiter{
		client:      client,
		maxRequests: maxRequests,
		window:      window,
	}
}

func (l *redisRateLimiter) Allow(ctx context.Context, clientID string) bool {
	key := redisPrefix + clientID
	now := time.Now().UnixMilli()
	windowStart := now - l.window.Milliseconds()

	pipe := l.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", formatMillis(windowStart))
	pipe.ZCard(ctx, key)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.Warn("redis rate limit indisponível, permitindo", "error", err)
		return true
	}

	count, ok := cmds[1].(*redis.IntCmd)
	if !ok {
		return true
	}

	if count.Val() >= int64(l.maxRequests) {
		return false
	}

	pipe = l.client.Pipeline()
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: formatMillis(now)})
	pipe.Expire(ctx, key, l.window)
	if _, err := pipe.Exec(ctx); err != nil {
		slog.Warn("redis rate limit add falhou, permitindo", "error", err)
		return true
	}

	return true
}

type inMemoryRateLimiter struct {
	mu          sync.Mutex
	entries     map[string]rateLimitEntry
	maxRequests int
	window      time.Duration
}

type rateLimitEntry struct {
	count     int
	startedAt time.Time
}

func newInMemoryRateLimiter(maxRequests int, window time.Duration) *inMemoryRateLimiter {
	return &inMemoryRateLimiter{
		entries:     make(map[string]rateLimitEntry),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (l *inMemoryRateLimiter) Allow(_ context.Context, clientID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry, exists := l.entries[clientID]
	if !exists || now.Sub(entry.startedAt) >= l.window {
		l.entries[clientID] = rateLimitEntry{count: 1, startedAt: now}
		return true
	}

	if entry.count >= l.maxRequests {
		return false
	}

	entry.count++
	l.entries[clientID] = entry
	return true
}

func formatMillis(ms int64) string {
	return time.UnixMilli(ms).UTC().Format(time.RFC3339Nano)
}
