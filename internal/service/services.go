package service

import (
	"github.com/tsoy/rental-rewards/internal/data"
	"github.com/tsoy/rental-rewards/internal/events"
)

type Services struct {
	Payment PaymentService
}

func NewServices(models data.Models, publisher events.Publisher) Services {
	return Services{
		Payment: NewPaymentService(models.Payments, publisher),
	}
}
