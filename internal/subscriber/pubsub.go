package subscriber

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/tsoy/rental-rewards/internal/service"
)

type Subscriber struct {
	client *pubsub.Client
	svc    *service.PaymentService
}

func NewSubscriber(client *pubsub.Client, svc *service.PaymentService) *Subscriber {
	return &Subscriber{client: client, svc: svc}
}

func (s *Subscriber) Register(ctx context.Context) {
	go s.startPaymentCompletedSubscriber(ctx)
}
