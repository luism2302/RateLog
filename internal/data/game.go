package data

import "time"

type Game struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Released        time.Time `json:"released"`
	BackgroundImage string    `json:"background_image"`
	Added           int       `json:"added"`
}

func NewGame(id int, name, released, image string) (Game, error) {
	date, err := time.Parse(time.DateOnly, released)
	if err != nil {
		return Game{}, err
	}

	return Game{ID: id, Name: name, Released: date, BackgroundImage: image}, nil
}

func SortGame(g1, g2 Game) int {
	return g2.Added - g1.Added
}
