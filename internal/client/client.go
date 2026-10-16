package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"unicode"

	"github.com/luism2302/RateLog/internal/data"
)

const (
	baseURL = "https://api.rawg.io/api"
)

type Client struct {
	apiKey string
	Client *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		Client: http.DefaultClient,
	}
}

func (c *Client) SearchGameByTitle(title string) ([]data.Game, error) {
	reqURL := fmt.Sprintf("%s/games?key=%s&search=%s&page_size=40", baseURL, c.apiKey, url.PathEscape(title))
	response, err := c.Client.Get(reqURL)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	games := []data.Game{}
	for _, result := range output.Results {
		if !strings.Contains(cleanTitle(result.Name), cleanTitle(title)) || result.Released == "" || result.BackgroundImage == "" {
			continue
		}

		game, err := data.NewGame(result.ID, result.Added, result.Name, result.Released, result.BackgroundImage)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	slices.SortFunc(games, data.SortGame)
	return games, nil
}

func readJSON(response *http.Response, data any) error {
	decoder := json.NewDecoder(response.Body)
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
