package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoRowsAffected = errors.New("No rows affected")

type UserGamesModel struct {
	DB *pgxpool.Pool
}

func (m UserGamesModel) Insert(userID, gameID int) error {
	query := `INSERT INTO users_games(user_id, game_id) VALUES ($1, $2) ON CONFLICT (user_id, game_id) DO NOTHING`
	args := []any{userID, gameID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)

	return err
}

func (m UserGamesModel) Delete(userID, gameID int) error {
	query := `DELETE FROM users_games WHERE user_id = $1 AND game_id = $2`
	args := []any{userID, gameID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tag, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return err
}
