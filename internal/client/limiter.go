package client

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	userRate      = 6 * time.Second
	userBurst     = 10
	userTimeLimit = 3 * time.Minute
	globalRate    = 500 * time.Millisecond
	globalBurst   = 5
	dailyRate     = 150 * time.Second
	dailyBurst    = 50
	defaultPause  = 10 * time.Minute
	minPause      = time.Minute
	maxPause      = time.Hour
)

var ErrRateLimited = errors.New("Rate limited")

type userLimiter struct {
	mu    sync.Mutex
	users map[int]*userClient
}

type globalLimiter struct {
	mu          sync.Mutex
	burst       *rate.Limiter
	daily       *rate.Limiter
	pausedUntil time.Time
}

type userClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func newUserLimiter() *userLimiter {
	l := &userLimiter{
		users: make(map[int]*userClient),
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			l.deleteInactive(time.Now())
		}
	}()
	return l
}

func newGlobalLimiter() *globalLimiter {
	daily := rate.NewLimiter(rate.Every(dailyRate), dailyBurst)
	daily.AllowN(time.Now(), dailyBurst-20)
	return &globalLimiter{
		burst: rate.NewLimiter(rate.Every(globalRate), globalBurst),
		daily: daily,
	}
}

func (u *userLimiter) reserve(quotaKey int, now time.Time) (*rate.Reservation, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, ok := u.users[quotaKey]
	if !ok {
		lim := rate.NewLimiter(rate.Every(userRate), userBurst)
		newUser := &userClient{limiter: lim}
		u.users[quotaKey] = newUser
		user = newUser
	}
	u.users[quotaKey].lastSeen = time.Now()

	r := user.limiter.ReserveN(now, 1)
	if !r.OK() || r.DelayFrom(now) > 0 {
		r.CancelAt(now)
		return nil, false
	}
	return r, true
}

func (u *userLimiter) deleteInactive(now time.Time) {
	u.mu.Lock()
	defer u.mu.Unlock()

	for key, user := range u.users {
		if now.Sub(user.lastSeen) > userTimeLimit {
			delete(u.users, key)
		}
	}
}

func (g *globalLimiter) allow(now time.Time) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if now.Before(g.pausedUntil) {
		return false
	}

	if !g.burst.AllowN(now, 1) {
		return false
	}

	if !g.daily.AllowN(now, 1) {
		return false
	}

	return true
}

func (g *globalLimiter) pauseFor(now time.Time, d time.Duration) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if now.Add(d).After(g.pausedUntil) {
		g.pausedUntil = now.Add(d)
		return true
	}

	return false
}

func (g *globalLimiter) retryAfter(header string, now time.Time) time.Duration {
	cleanHeader := strings.TrimSpace(header)
	if cleanHeader == "" {
		return defaultPause
	}

	if secs, err := strconv.Atoi(cleanHeader); err == nil && secs >= 0 {
		if secs > int(maxPause/time.Second) {
			return maxPause
		}
		return clamp(time.Duration(secs) * time.Second)
	}

	if t, err := http.ParseTime(cleanHeader); err == nil {
		return clamp(t.Sub(now))
	}

	return defaultPause
}

func clamp(d time.Duration) time.Duration {
	if d < minPause {
		return minPause
	}
	if d > maxPause {
		return maxPause
	}

	return d
}
