package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users    UserModel
	Payments PaymentModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users:    UserModel{DB: db},
		Payments: PaymentModel{DB: db},
	}
}
