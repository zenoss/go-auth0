package http

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var (
	limiters              = make(map[string]*rate.Limiter)
	muLimiters            = sync.Mutex{}
	cfgLimit   rate.Limit = 2
	cfgBurst   int        = 5
)

func SetRateLimitConfig(limit rate.Limit) {
	muLimiters.Lock()
	defer muLimiters.Unlock()
	cfgLimit = limit
	for _, l := range limiters {
		l.SetLimit(limit)
	}
}

func SetBurstRateConfig(burst int) {
	muLimiters.Lock()
	defer muLimiters.Unlock()
	cfgBurst = burst
	for _, l := range limiters {
		l.SetBurst(burst)
	}
}

func getLimiter(key string) *rate.Limiter {
	muLimiters.Lock()
	defer muLimiters.Unlock()

	limiter, ok := limiters[key]
	if !ok {
		limiter = rate.NewLimiter(cfgLimit, cfgBurst)
		limiters[key] = limiter
	}
	return limiter
}

func GetRequestLimiter(req *http.Request) *rate.Limiter {
	return getLimiter(fmt.Sprintf("%s %s", req.Method, req.URL.Path))
}
