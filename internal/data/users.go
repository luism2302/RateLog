package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordHashNil = errors.New("password hash is nil")
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrRecordNotFound  = errors.New("record not found")
)

type User struct {
	ID        int
	Email     string
	Name      string
	Password  password
	CreatedAt time.Time
	Version   int
}

func NewUser(email, name, plaintextPassword string) (User, error) {
	user := User{Email: email, Name: name}
	if err := user.Password.Set(plaintextPassword); err != nil {
		return User{}, err
	}

	if user.Password.hash == nil {
		return User{}, ErrPasswordHashNil
	}

	return user, nil
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m UserModel) Insert(user User) error {
	query := `INSERT INTO users(name, email, password_hash) VALUES ($1, $2, $3)`
	args := []any{user.Name, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, args...)
	if err != nil {
		var pgxErr *pgconn.PgError
		switch {
		case errors.As(err, &pgxErr) && pgxErr.Code == "23505" && pgxErr.ConstraintName == "users_email_key":
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetUserByEmail(email string) (User, error) {
	query := `SELECT id, name, email, password_hash, created_at, version FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	if err := m.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password.hash, &user.CreatedAt, &user.Version); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return User{}, ErrRecordNotFound
		default:
			return User{}, err
		}
	}

	return user, nil
}

func (m UserModel) GetUserForToken(tokenPlain string) (User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlain))

	query := `SELECT users.id, users.email, users.name, users.created_at, users.password_hash, users.version FROM users JOIN tokens ON tokens.user_id = users.id WHERE tokens.hash = $1 AND tokens.expires > $2`
	args := []any{tokenHash[:], time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	if err := m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.Password.hash, &user.Version); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return User{}, ErrRecordNotFound
		default:
			return User{}, err
		}
	}

	return user, nil
}
