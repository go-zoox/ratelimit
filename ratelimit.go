package ratelimit

import (
	"time"

	"github.com/go-zoox/counter"
	"github.com/go-zoox/counter/bucket"
)

// RateLimit is a rate limiter, support custom bucket (memory, redis, other databases).
type RateLimit struct {
	counter.Counter
	maxCount int64
	//
	maxAge time.Duration
}

// Status is the status of the rate limit.
type Status struct {
	Total     int64
	Remaining int64
	// the time when the rate limit will be reset, unit: milliseconds
	ResetAfter  int64
	IsExeceeded bool
}

// New creates a new rate limiter.
func New(bucket bucket.Bucket, namespace string, maxAge time.Duration, maxCount int64) *RateLimit {
	return &RateLimit{
		Counter:  *counter.New(bucket, namespace, maxAge),
		maxCount: maxCount,
		maxAge:   maxAge,
	}
}

// Total is the total number of occurrences during a period.
func (r *RateLimit) Total(id string) int64 {
	return r.maxCount
}

// Remaining is the number of occurrences remaining during a period.
func (r *RateLimit) Remaining(id string) int64 {
	count, err := r.Count(id)
	if err != nil {
		// panic(err)
		return r.maxCount
	}

	return r.maxCount - count.Count
}

// ResetAfter is the time when the rate limit will be reset.
func (r *RateLimit) ResetAfter(id string) int64 {
	count, err := r.Count(id)
	if err != nil {
		panic(err)
	}

	return count.ExpiresAt - time.Now().UnixMilli()
}

// ResetAt is the time when the rate limit will be reset.
func (r *RateLimit) ResetAt(id string) int64 {
	count, err := r.Count(id)
	if err != nil {
		panic(err)
	}

	return count.ExpiresAt
}

// IsExceeded is the rate limit is exceeded.
func (r *RateLimit) IsExceeded(id string) bool {
	return r.Remaining(id) < 0
}

// Status is the status of the rate limit.
func (r *RateLimit) Status(id string) (*Status, error) {
	count, err := r.Count(id)
	if err != nil {
		return nil, err
	}

	if count.ExpiresAt == 0 {
		return &Status{
			Total:       r.maxCount,
			Remaining:   r.maxCount,
			ResetAfter:  int64(r.maxAge / time.Millisecond),
			IsExeceeded: false,
		}, nil
	}

	maxCount := r.maxCount
	remaining := maxCount - count.Count
	resetAfter := count.ExpiresAt - time.Now().UnixMilli()

	return &Status{
		Total:       maxCount,
		Remaining:   remaining,
		ResetAfter:  resetAfter,
		IsExeceeded: remaining < 0,
	}, nil
}

// NewMemory creates a in-memory ratelimit.
func NewMemory(namespace string, maxAge time.Duration, maxCount int64) *RateLimit {
	return New(
		bucket.NewMemory(),
		namespace,
		maxAge,
		maxCount,
	)
}

// NewRedis creates a redis-baed ratelimit.
func NewRedis(namespace string, maxAge time.Duration, maxCount int64, cfg *bucket.RedisConfig) (*RateLimit, error) {
	b, err := bucket.NewRedis(cfg)
	if err != nil {
		return nil, err
	}

	return New(b, namespace, maxAge, maxCount), nil
}
