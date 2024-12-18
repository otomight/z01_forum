package servsetup

import (
	"sync"
	"time"
)

type RateLimiter struct {
	Visitors map[string]*Visitor
	Mu		sync.Mutex
	Rate	int
	Window	time.Duration
}

type Visitor struct {
	LastRequest	time.Time
	Requests	int
}


func newRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Visitors:	make(map[string]*Visitor),
		Rate:		rate,
		Window:		window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	var	visitor	*Visitor
	var	exists	bool

	rl.Mu.Lock()
	defer rl.Mu.Unlock()
	visitor, exists = rl.Visitors[ip]
	if !exists || time.Since(visitor.LastRequest) > rl.Window {
		rl.Visitors[ip] = &Visitor{
			LastRequest:	time.Now(),
			Requests:		1,
		}
		return true
	}
	if visitor.Requests > rl.Rate {
		return false
	}
	visitor.Requests++
	visitor.LastRequest = time.Now()
	return true
}
