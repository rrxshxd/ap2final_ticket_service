package payment

type PaymentRequest struct {
	Amount   float64
	Currency string
}

type PaymentResponse struct {
	PaymentID string
}
