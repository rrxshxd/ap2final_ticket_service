package models

import "time"

type TicketStatus string

const (
	TicketStatusReserved  TicketStatus = "RESERVED"
	TicketStatusPaid      TicketStatus = "PAID"
	TicketStatusCancelled TicketStatus = "CANCELLED"
)

type Ticket struct {
	ID            string       `bson:"_id"`
	SessionID     string       `bson:"session_id"`
	MovieID       string       `bson:"movie_id"`
	SeatNumber    string       `bson:"seat_number"`
	Price         float64      `bson:"price"`
	Status        TicketStatus `bson:"status"`
	UserID        string       `bson:"user_id"`
	PurchaseDate  time.Time    `bson:"purchase_date"`
	PaymentMethod string       `bson:"payment_method"`
}
