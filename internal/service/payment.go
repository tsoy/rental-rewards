package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/tsoy/rental-rewards/internal/data"
	"github.com/tsoy/rental-rewards/internal/events"
	"time"
)

type PaymentService struct {
	Payments  data.PaymentModel
	Publisher events.Publisher
}

func NewPaymentService(payments data.PaymentModel, publisher events.Publisher) PaymentService {
	return PaymentService{
		Payments:  payments,
		Publisher: publisher,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment *data.Payment) error {
	payment.ID = uuid.New()
	payment.Status = data.StatusCompleted
	now := time.Now().UTC()
	payment.CompletedAt = &now
	if payment.Currency == "" {
		payment.Currency = data.USD
	}

	// Insert payment into DB
	if err := s.Payments.Insert(ctx, payment); err != nil {
		return err
	}

	// Build and publish event
	evt := events.PaymentCompletedEvent{
		EventType:    "payment.completed",
		EventVersion: "1.0",
		ID:           uuid.NewString(),
		OccurredAt:   time.Now(),
		Payment:      events.FromDBPayment(*payment),
		Metadata:     events.Metadata{Source: "api-service"},
	}

	if err := s.Publisher.PublishPaymentCompleted(ctx, evt); err != nil {
		// Log error; maybe retry; decide strategy
		return err
	}

	return nil
}

func (s *PaymentService) HandlePaymentCompleted(ctx context.Context, event events.PaymentCompletedEvent) error {
	return nil
}
