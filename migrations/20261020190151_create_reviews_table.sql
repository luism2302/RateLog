-- +goose Up
CREATE TABLE IF NOT EXISTS reviews(
     id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
     user_id BIGINT NOT NULL,
     game_id BIGINT NOT NULL,
     body TEXT,
     rating SMALLINT NOT NULL CHECK(rating BETWEEN 1 AND 5),
     created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
     UNIQUE (user_id, game_id),
     FOREIGN KEY (user_id, game_id) REFERENCES users_games (user_id, game_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS reviews;
