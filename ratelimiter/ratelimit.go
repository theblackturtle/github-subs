package ratelimiter

import (
	"sync"
	"time"
)

type RateLimit struct {
	// Rate limit enforcement fields
	rateLimit time.Duration
	lastLock  sync.Mutex
	last      time.Time
}

func NewRateLimiter() *RateLimit {
	r := &RateLimit{
		last: time.Now().Truncate(10 * time.Minute),
	}
	return r
}

// SetRateLimit sets the minimum wait between checks.
func (r *RateLimit) SetRateLimit(min time.Duration) {
	r.rateLimit = min
}

// CheckRateLimit blocks until the minimum wait since the last call.
func (r *RateLimit) CheckRateLimit() {
	if r.rateLimit == time.Duration(0) {
		return
	}

	r.lastLock.Lock()
	defer r.lastLock.Unlock()

	if delta := time.Now().Sub(r.last); r.rateLimit > delta {
		time.Sleep(delta)
	}
	r.last = time.Now()
}
