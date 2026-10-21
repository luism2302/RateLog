package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Review struct {
	UserID    int
	GameID    int
	Body      string
	Rating    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GameReview struct {
	Game   Game
	Review Review
}

func NewGameReview(game Game, review Review) GameReview {
	return GameReview{
		Game:   game,
		Review: review,
	}
}

type ReviewsModel struct {
	DB *pgxpool.Pool
}

func NewReview(userID, gameID int, body string, rating int) Review {
	return Review{
		UserID: userID,
		GameID: gameID,
		Body:   body,
		Rating: rating,
	}
}

func (m ReviewsModel) Insert(review Review) error {
	query := `INSERT INTO reviews(user_id, game_id, body, rating) VALUES($1, $2, $3, $4)`
	args := []any{review.UserID, review.GameID, review.Body, review.Rating}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := m.DB.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (m ReviewsModel) GetReviewsForUser(userID int) ([]GameReview, error) {
	query := `SELECT games.title, games.background_image_link, reviews.body, reviews.rating, reviews.updated_at FROM reviews JOIN games ON games.id = reviews.game_id WHERE user_id = $1 ORDER BY reviews.updated_at DESC`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return []GameReview{}, ErrRecordNotFound
		default:
			return []GameReview{}, err
		}
	}

	reviews := []GameReview{}
	for rows.Next() {
		var review GameReview
		if err := rows.Scan(&review.Game.Name, &review.Game.BackgroundImage, &review.Review.Body, &review.Review.Rating, &review.Review.UpdatedAt); err != nil {
			return []GameReview{}, err
		}

		reviews = append(reviews, review)
	}

	if rows.Err() != nil {
		return []GameReview{}, rows.Err()
	}

	return reviews, err
}
