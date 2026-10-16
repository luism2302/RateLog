package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
