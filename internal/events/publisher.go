package events

import "context"

type Publisher interface {
	PublishPaymentCompleted(ctx context.Context, evt PaymentCompletedEvent) error
}
