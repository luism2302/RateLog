package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Game struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Released        time.Time `json:"released"`
	BackgroundImage string    `json:"background_image"`
	Added           int       `json:"added"`
}

func NewGame(id, added int, name, released, image string) (Game, error) {
	date, err := time.Parse(time.DateOnly, released)
	if err != nil {
		return Game{}, err
	}

	return Game{ID: id, Name: name, Released: date, BackgroundImage: image, Added: added}, nil
}

func SortGame(g1, g2 Game) int {
	return g2.Added - g1.Added
}

type GameModel struct {
	DB *pgxpool.Pool
}

func (m GameModel) Insert(game Game) error {
	query := `INSERT INTO games(id, title, release_date, background_image_link) VALUES($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`
	args := []any{game.ID, game.Name, game.Released, game.BackgroundImage}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)

	return err
}
