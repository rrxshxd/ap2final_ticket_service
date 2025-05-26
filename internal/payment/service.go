package payment

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ProcessPayment(ctx context.Context, amount float64, currency string) (string, error)
}

type mockPaymentService struct{}

func NewMockPaymentService() Service {
	return &mockPaymentService{}
}

func (s *mockPaymentService) ProcessPayment(ctx context.Context, amount float64, currency string) (string, error) {
	objectID := primitive.NewObjectID()
	return "mock-pay-" + objectID.Hex(), nil
}
