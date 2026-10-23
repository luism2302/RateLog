package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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

func (m GameModel) GetGamesForUser(userID int) ([]Game, error) {
	query := `SELECT games.id, games.title, games.release_date, games.background_image_link
	FROM games 
	JOIN users_games ON games.id = users_games.game_id AND users_games.user_id = $1 ORDER BY users_games.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, userID)
	if err != nil {
		return []Game{}, err
	}
	defer rows.Close()

	games := []Game{}
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Name, &game.Released, &game.BackgroundImage); err != nil {
			return []Game{}, err
		}
		games = append(games, game)
	}

	if rows.Err() != nil {
		return []Game{}, rows.Err()
	}

	return games, nil
}

func (m GameModel) GetGameByID(gameID int) (Game, error) {
	query := `SELECT id, title, release_date, background_image_link FROM games WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var game Game
	if err := m.DB.QueryRow(ctx, query, gameID).Scan(&game.ID, &game.Name, &game.Released, &game.BackgroundImage); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return Game{}, ErrRecordNotFound
		default:
			return Game{}, err
		}
	}

	return game, nil
}
