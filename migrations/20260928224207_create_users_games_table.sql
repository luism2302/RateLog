-- +goose Up
CREATE TABLE IF NOT EXISTS users_games(
     user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
     game_id BIGINT NOT NULL REFERENCES games ON DELETE CASCADE,
     PRIMARY KEY (user_id, game_id)
);

-- +goose Down
DROP TABLE IF EXISTS users_games;
