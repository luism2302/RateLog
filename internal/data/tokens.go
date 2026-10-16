package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luism2302/RateLog/internal/validator"
)

type Token struct {
	UserID    int
	Plaintext string
	Hash      []byte
	Expires   time.Time
}

func generateToken(userID int, expires time.Duration) Token {
	token := Token{
		UserID:    userID,
		Plaintext: rand.Text(),
		Expires:   time.Now().Add(expires),
	}

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token
}

func ValidateToken(token string, v *validator.Validator) {
	v.Check(token != "", "token", "Cannot be blank")
	v.Check(len(token) == 26, "token", "Must be 26 bytes long")
}

type TokenModel struct {
	DB *pgxpool.Pool
}

func (m TokenModel) insert(token Token) error {
	query := `INSERT INTO tokens(hash, user_id, expires) VALUES($1, $2, $3)`
	args := []any{token.Hash, token.UserID, token.Expires}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)

	return err
}

func (m TokenModel) NewToken(userID int, expires time.Duration) (Token, error) {
	token := generateToken(userID, expires)
	err := m.insert(token)

	return token, err
}

func (m TokenModel) Delete(tokenPlain string) error {
	hash := sha256.Sum256([]byte(tokenPlain))

	query := `DELETE FROM tokens WHERE hash = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, hash[:])

	return err
}
