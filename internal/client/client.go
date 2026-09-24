package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

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
	reqURL := fmt.Sprintf("%s/games?key=%s&search=%s&exclude_additions=true&page_size=40", baseURL, c.apiKey, url.PathEscape(title))
	fmt.Println(reqURL)
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
		if !strings.Contains(strings.ToLower(result.Name), strings.ToLower(title)) || result.Released == "" {
			continue
		}

		game, err := data.NewGame(result.ID, result.Name, result.Released, result.BackgroundImage)
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
