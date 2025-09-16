package data

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	StatusCreated   = "created"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

const (
	USD = "USD"
	CAD = "CAD"
	CNY = "CNY"
)

var Currencies = map[string]struct{}{
	USD: {},
	CAD: {},
	CNY: {},
}

type Payment struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	AmountCents int64      `json:"amount_cents"`
	Currency    string     `json:"currency"`
	Status      string     `json:"status"`
	ExternalRef *string    `json:"external_ref,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

var ErrDuplicateTransaction = errors.New("payment with external_ref already exists")
var ErrInvalidUserId = errors.New("invalid user ID")

type PaymentModel struct {
	DB *pgxpool.Pool
}

func (p PaymentModel) Insert(ctx context.Context, payment *Payment) error {
	//payment.ID = uuid.New()
	//payment.Status = StatusCompleted
	//now := time.Now().UTC()
	//payment.CompletedAt = &now
	//if payment.Currency == "" {
	//	payment.Currency = USD
	//}
	tx, err := p.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	//var ext *string
	//if payment.ExternalRef != nil && *payment.ExternalRef != "" {
	//
	//}
	err = p.DB.QueryRow(
		ctx,
		`INSERT INTO payments (id,user_id,amount_cents,currency,status,external_ref,completed_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at`,
		payment.ID, payment.UserID, payment.AmountCents, payment.Currency, payment.Status, payment.ExternalRef, payment.CompletedAt,
	).Scan(&payment.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" && pgErr.ConstraintName == "payments_external_ref_uidx":
				return ErrDuplicateTransaction
			case pgErr.Code == "23503" && pgErr.ConstraintName == "payments_user_id_fkey":
				return ErrInvalidUserId
			}
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p PaymentModel) Get(ctx context.Context, paymentId uuid.UUID) (*Payment, error) {
	pm := &Payment{}
	err := p.DB.QueryRow(ctx, `SELECT id,user_id,amount_cents,currency,status,external_ref,created_at,completed_at 
		FROM payments WHERE id=$1`, paymentId).Scan(
		&pm.ID, &pm.UserID, &pm.AmountCents,
		&pm.Currency, &pm.Status, &pm.ExternalRef, &pm.CreatedAt, &pm.CompletedAt)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return pm, nil
}
