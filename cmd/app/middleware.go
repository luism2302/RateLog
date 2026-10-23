package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/internal/validator"
	"golang.org/x/time/rate"
)

const (
	looseRate   = 5
	strictRate  = 2
	looseBurst  = 10
	strictBurst = 4
)

func (app *application) authenticate(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		authCookie, err := request.Cookie(authCookieName)
		if err != nil {
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		token := authCookie.Value
		v := validator.New()

		if data.ValidateToken(token, v); !v.Valid() {
			removeAuthCookie(responseWriter)
			http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
			return
		}

		user, err := app.models.Users.GetUserForToken(token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				removeAuthCookie(responseWriter)
				http.Redirect(responseWriter, request, "/login", http.StatusSeeOther)
				return
			default:
				app.serverErrorResponse(responseWriter, request, err)
				return
			}
		}
		request = app.setAuthenticatedUser(request, user)
		next.ServeHTTP(responseWriter, request)
	})
}

func (app *application) rateLimit(next http.Handler) http.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for key, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, "/static/") {
			next.ServeHTTP(responseWriter, request)
			return
		}

		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			ip = request.RemoteAddr
		}
		endpoint := request.URL.Path
		rateCap := looseRate
		burst := looseBurst
		segment := "normal"

		if request.Method == http.MethodPost && (endpoint == "/login" || endpoint == "/signup") {
			rateCap = strictRate
			burst = strictBurst
			segment = "strict"
		}

		key := fmt.Sprintf("%s:%s", ip, segment)

		mu.Lock()
		if _, found := clients[key]; !found {
			clients[key] = &client{limiter: rate.NewLimiter(rate.Limit(rateCap), burst)}
		}

		clients[key].lastSeen = time.Now()

		if !clients[key].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(responseWriter, request)
			return
		}

		mu.Unlock()
		next.ServeHTTP(responseWriter, request)
	})
}

func (app *application) commonHeaders(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Add("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(responseWriter, request)
	})
}
