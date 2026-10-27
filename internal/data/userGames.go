package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	if _, err := m.DB.Exec(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23503":
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
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

func (m UserGamesModel) HasGame(userID, gameID int) error {
	query := `SELECT user_id, game_id  FROM users_games WHERE user_id =  $1 AND game_id = $2`
	args := []any{userID, gameID}

	var output struct {
		userID int
		gameID int
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRow(ctx, query, args...).Scan(&output.userID, &output.gameID); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}
