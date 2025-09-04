package data

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  *string   `json:"full_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (u UserModel) Insert(ctx context.Context, user *User) error {
	//id := uuid.New()
	user.ID = uuid.New()
	err := u.DB.QueryRow(
		ctx,
		`INSERT INTO users (id,email,full_name) VALUES ($1,$2,$3) RETURNING created_at`,
		user.ID, user.Email, user.FullName,
	).Scan(&user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrDuplicateEmail
			}
		}
		return err
		//switch {
		//case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
		//	return ErrDuplicateEmail
		//default:
		//	return err
		//}
	}
	//if err != nil {
	//	return err
	//}
	return nil
}
