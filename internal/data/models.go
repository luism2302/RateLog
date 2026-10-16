package data

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Users     UserModel
	Tokens    TokenModel
	Games     GameModel
	UserGames UserGamesModel
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Users:     UserModel{DB: pool},
		Tokens:    TokenModel{DB: pool},
		Games:     GameModel{DB: pool},
		UserGames: UserGamesModel{DB: pool},
	}
}
