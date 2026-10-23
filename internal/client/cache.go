package client

import (
	"slices"
	"sync"
	"time"

	"github.com/luism2302/RateLog/internal/data"
)

const (
	maxEntries         = 200
	maxMissingEntries  = 2000
	timeInCache        = 5 * time.Minute
	missingTimeInCache = time.Hour
)

type idPayload struct {
	game      data.Game
	expiresAt time.Time
}
type gameIDCache struct {
	mu      sync.Mutex
	entries map[int]idPayload
	missing map[int]time.Time
}

type titlePayload struct {
	games     []data.Game
	expiresAt time.Time
}

type gameTitleCache struct {
	mu      sync.Mutex
	entries map[string]titlePayload
}

type cache struct {
	idCache    *gameIDCache
	titleCache *gameTitleCache
}

func newGameIDCache() *gameIDCache {
	return &gameIDCache{
		entries: make(map[int]idPayload),
		missing: make(map[int]time.Time),
	}
}

func newGameTitleCache() *gameTitleCache {
	return &gameTitleCache{
		entries: make(map[string]titlePayload),
	}
}

func newCache() cache {
	cache := cache{
		idCache:    newGameIDCache(),
		titleCache: newGameTitleCache(),
	}
	go cache.cleanCache(timeInCache)
	return cache
}

func (c cache) cleanCache(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.idCache.deleteExpired()
		c.titleCache.deleteExpired()
	}
}

func (c *gameTitleCache) get(title string) ([]data.Game, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.entries[title]
	if !ok {
		return nil, false
	}

	if value.expiresAt.Before(time.Now()) {
		return nil, false
	}

	return slices.Clone(value.games), true
}

func (c *gameTitleCache) set(title string, games []data.Game) {
	c.mu.Lock()
	defer c.mu.Unlock()

	payload := titlePayload{
		games:     slices.Clone(games),
		expiresAt: time.Now().Add(timeInCache),
	}

	if _, ok := c.entries[title]; ok {
		c.entries[title] = payload
		return
	}

	if len(c.entries) >= maxEntries {
		for key, payload := range c.entries {
			if payload.expiresAt.Before(time.Now()) {
				delete(c.entries, key)
			}
		}
	}

	// still full
	if len(c.entries) >= maxEntries {
		oldestKey := ""
		var oldestTime time.Time
		i := 0
		for key, payload := range c.entries {
			if i == 0 {
				oldestKey = key
				oldestTime = payload.expiresAt
				i++
				continue
			}
			if payload.expiresAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = payload.expiresAt
			}
		}
		delete(c.entries, oldestKey)
	}

	c.entries[title] = payload
}

func (c *gameTitleCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, val := range c.entries {
		if val.expiresAt.Before(time.Now()) {
			delete(c.entries, key)
		}
	}
}

func (c *gameIDCache) get(id int) (data.Game, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.entries[id]
	if !ok {
		return data.Game{}, false
	}
	if value.expiresAt.Before(time.Now()) {
		return data.Game{}, false
	}
	return value.game, true
}

func (c *gameIDCache) set(id int, game data.Game) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.missing, id)

	payload := idPayload{
		game:      game,
		expiresAt: time.Now().Add(timeInCache),
	}

	if _, ok := c.entries[id]; ok {
		c.entries[id] = payload
		return
	}

	if len(c.entries) >= maxEntries {
		for key, payload := range c.entries {
			if payload.expiresAt.Before(time.Now()) {
				delete(c.entries, key)
			}
		}
	}

	// still full
	if len(c.entries) >= maxEntries {
		oldestKey := 0
		var oldestTime time.Time
		i := 0
		for key, payload := range c.entries {
			if i == 0 {
				oldestKey = key
				oldestTime = payload.expiresAt
				i++
				continue
			}
			if payload.expiresAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = payload.expiresAt
			}
		}
		delete(c.entries, oldestKey)
	}

	c.entries[id] = payload
}

func (c *gameIDCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, val := range c.entries {
		if val.expiresAt.Before(time.Now()) {
			delete(c.entries, key)
		}
	}
	for key, val := range c.missing {
		if val.Before(time.Now()) {
			delete(c.missing, key)
		}
	}
}

func (c *gameIDCache) markMissing(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.missing[id]; ok {
		c.missing[id] = time.Now().Add(missingTimeInCache)
		return
	}

	if len(c.missing) >= maxMissingEntries {
		for key, t := range c.missing {
			if t.Before(time.Now()) {
				delete(c.missing, key)
			}
		}
	}

	// still full
	if len(c.missing) >= maxMissingEntries {
		oldestKey := 0
		var oldestTime time.Time
		i := 0
		for key, t := range c.missing {
			if i == 0 {
				oldestKey = key
				oldestTime = t
				i++
				continue
			}
			if t.Before(oldestTime) {
				oldestKey = key
				oldestTime = t
			}
		}
		delete(c.missing, oldestKey)
	}

	c.missing[id] = time.Now().Add(missingTimeInCache)
}

func (c *gameIDCache) getMissing(id int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.missing[id]
	if !ok {
		return false
	}

	if val.Before(time.Now()) {
		delete(c.missing, id)
		return false
	}

	return true
}
