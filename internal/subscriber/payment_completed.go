package subscriber

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/tsoy/rental-rewards/internal/events"
	"log"
)

func (s *Subscriber) startPaymentCompletedSubscriber(ctx context.Context) {
	sub := s.client.Subscription("rewards-worker-sub")

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		defer msg.Ack()

		var event events.PaymentCompletedEvent

		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("failed to unmarshal payment.compoleted event: %v", err)
			return
		}

		if err := s.svc.HandlePaymentCompleted(ctx, event); err != nil {
			log.Printf("failed to handle payment.completed: %v", err)
			return
		}

		log.Printf("processed payment.completed for payment: %v", event.Payment.ID)
	})

	if err != nil {
		log.Fatalf("subscriber for payment.completed stopped: %v", err)
	}
	//go s.startPaymentCompletedSubscriber(ctx)
}
