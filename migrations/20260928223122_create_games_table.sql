-- +goose Up
CREATE TABLE IF NOT EXISTS games(
     id BIGINT PRIMARY KEY,
     title TEXT NOT NULL UNIQUE,
     release_date TIMESTAMP(0) WITH TIME ZONE NOT NULL,
     background_image_link TEXT NOT NULL,
     created_at TIMESTAMP (0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS games;
