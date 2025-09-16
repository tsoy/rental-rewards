package events

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
)

type PubSubPublisher struct {
	client *pubsub.Client
}

// NewPubSubPublisher creates a Pub/Sub-backed Publisher.
func NewPubSubPublisher(client *pubsub.Client) *PubSubPublisher {
	return &PubSubPublisher{client: client}
}

func (p *PubSubPublisher) PublishPaymentCompleted(ctx context.Context, evt PaymentCompletedEvent) error {
	topic := p.client.Topic("payment.completed")

	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	res := topic.Publish(ctx, &pubsub.Message{Data: data})
	_, err = res.Get(ctx)
	return err
}
