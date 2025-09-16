package events

import (
	"github.com/google/uuid"
	"github.com/tsoy/rental-rewards/internal/data"
	"time"
)

type Payment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	AmountCents int64     `json:"amount_cents"`
	Currency    string    `json:"currency"`
	CompletedAt time.Time `json:"completed_at"`
}

type PaymentCompletedEvent struct {
	EventType    string    `json:"event_type"`
	EventVersion string    `json:"event_version"`
	ID           string    `json:"id"`
	OccurredAt   time.Time `json:"occurred_at"`
	Payment      Payment   `json:"payment"`
	Metadata     Metadata  `json:"metadata"`
}

type Metadata struct {
	TraceID string `json:"trace_id"`
	Source  string `json:"source"`
}

func FromDBPayment(p data.Payment) Payment {
	return Payment{
		ID:          p.ID,
		UserID:      p.UserID,
		AmountCents: p.AmountCents,
		Currency:    p.Currency,
		CompletedAt: *p.CompletedAt,
	}
}
