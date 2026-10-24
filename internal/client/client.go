package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/luism2302/RateLog/internal/data"
)

var (
	ErrNotFound    = errors.New("Not Found")
	ErrUnavailable = errors.New("Unavailable")
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.rawg.io",
}

type Client struct {
	apiKey       string
	client       *http.Client
	cache        cache
	userLimits   *userLimiter
	globalLimits *globalLimiter
}

func New(apiKey string) *Client {
	httpClient := &http.Client{
		Timeout: 8 * time.Second,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Client{
		apiKey:       apiKey,
		client:       httpClient,
		cache:        newCache(),
		userLimits:   newUserLimiter(),
		globalLimits: newGlobalLimiter(),
	}
}

type params struct {
	calledFrom string
	id         int
	title      string
}

func (c *Client) SearchGameByID(quotaKey, id int, ctx context.Context) (data.Game, error) {
	if quotaKey <= 0 {
		return data.Game{}, errors.New("invalid user id for api call")
	}
	if game, ok := c.cache.idCache.get(id); ok {
		return game, nil
	}
	if c.cache.idCache.getMissing(id) {
		return data.Game{}, ErrNotFound
	}

	response, err := c.makeAPIRequest(quotaKey, ctx, params{calledFrom: "id", id: id})
	if err != nil {
		return data.Game{}, err
	}

	defer response.Body.Close()

	var output struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		Released        string `json:"released"`
		BackgroundImage string `json:"background_image"`
		Added           int    `json:"added"`
	}
	if err := readJSON(response, &output); err != nil {
		return data.Game{}, fmt.Errorf("failed parsing json: %w", ErrUnavailable)
	}

	if output.ID != id || strings.TrimSpace(output.Name) == "" {
		c.cache.idCache.markMissing(id)
		return data.Game{}, ErrNotFound

	}

	if !isSafeImage(output.BackgroundImage) {
		c.cache.idCache.markMissing(id)
		return data.Game{}, ErrNotFound
	}

	game, err := data.NewGame(output.ID, output.Added, output.Name, output.Released, output.BackgroundImage)
	if err != nil {
		c.cache.idCache.markMissing(id)
		return data.Game{}, ErrNotFound
	}

	c.cache.idCache.set(id, game)
	return game, nil
}

func (c *Client) SearchGameByTitle(quotaKey int, title string, ctx context.Context) ([]data.Game, error) {
	if quotaKey <= 0 {
		return nil, errors.New("invalid user id for api call")
	}
	normalizedTitle := strings.ToLower(strings.TrimSpace(title))
	if games, ok := c.cache.titleCache.get(normalizedTitle); ok {
		return games, nil
	}

	response, err := c.makeAPIRequest(quotaKey, ctx, params{calledFrom: "title", title: title})
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var output struct {
		Results []struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Released        string `json:"released"`
			BackgroundImage string `json:"background_image"`
			Added           int    `json:"added"`
		} `json:"results"`
	}

	if err := readJSON(response, &output); err != nil {
		return nil, fmt.Errorf("failed parsing json: %w", ErrUnavailable)
	}

	games := []data.Game{}
	for _, result := range output.Results {
		if !strings.Contains(cleanTitle(result.Name), cleanTitle(title)) || result.Released == "" || result.BackgroundImage == "" {
			continue
		}
		if !isSafeImage(result.BackgroundImage) {
			continue
		}

		game, err := data.NewGame(result.ID, result.Added, result.Name, result.Released, result.BackgroundImage)
		if err != nil {
			continue
		}

		games = append(games, game)
		c.cache.idCache.set(game.ID, game)
	}

	slices.SortFunc(games, data.SortGame)
	c.cache.titleCache.set(normalizedTitle, games)
	return games, nil
}

func readJSON(response *http.Response, data any) error {
	body := io.LimitReader(response.Body, 2<<20)
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func cleanTitle(title string) string {
	var b strings.Builder
	for _, char := range title {
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) {
			b.WriteRune(char)
		}
	}

	return strings.ToLower(b.String())
}

func isSafeImage(link string) bool {
	imageURL, err := url.Parse(link)
	if err != nil {
		return false
	}
	if imageURL.User != nil || imageURL.Port() != "" || imageURL.Scheme != "https" || imageURL.Host != "media.rawg.io" {
		return false
	}

	return true
}

func classifyStatusCode(status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("rawg status: %d. err: %w", status, ErrNotFound)
	case http.StatusTooManyRequests:
		return fmt.Errorf("rawg status: %d. err: %w", status, ErrUnavailable)
	default:
		return fmt.Errorf("rawg status: %d. err: %w", status, ErrUnavailable)
	}
}

func (c *Client) makeAPIRequest(quotaKey int, ctx context.Context, values params) (*http.Response, error) {
	reqURL := baseURL
	urlValues := url.Values{}
	urlValues.Set("key", c.apiKey)

	switch values.calledFrom {
	case "id":
		reqURL.Path = fmt.Sprintf("/api/games/%d", values.id)
	case "title":
		reqURL.Path = "/api/games"
		urlValues.Set("search", values.title)
		urlValues.Set("page_size", "40")
	default:
		return nil, errors.New("unknown request")

	}

	if !c.userLimits.allow(quotaKey) {
		return nil, ErrRateLimited
	}

	now := time.Now()
	if !c.globalLimits.allow(now) {
		return nil, fmt.Errorf("rawg: global limit or pause: %w", ErrUnavailable)
	}

	reqURL.RawQuery = urlValues.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("rawg: request failed: %w", ErrUnavailable)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rawg: request failed: %w", ErrUnavailable)
	}

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfter := response.Header.Get("Retry-After")
		c.globalLimits.pauseFor(now, c.globalLimits.retryAfter(retryAfter, now))
	}

	err = classifyStatusCode(response.StatusCode)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			response.Body.Close()
			if values.calledFrom == "id" {
				c.cache.idCache.markMissing(values.id)
				return nil, ErrNotFound
			} else if values.calledFrom == "title" {
				return nil, fmt.Errorf("rawg search: status 404: %w", ErrUnavailable)
			}
			return nil, err
		default:
			response.Body.Close()
			return nil, err
		}
	}

	return response, nil
}
