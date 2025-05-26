package models

import "time"

type TicketStatus string

const (
	TicketStatusReserved  TicketStatus = "RESERVED"
	TicketStatusPaid      TicketStatus = "PAID"
	TicketStatusCancelled TicketStatus = "CANCELLED"
)

type Ticket struct {
	ID            string       `bson:"-"`
	SessionID     string       `bson:"-"`
	MovieID       string       `bson:"-"`
	SeatNumber    string       `bson:"-"`
	Price         float64      `bson:"-"`
	Status        TicketStatus `bson:"-"`
	UserID        string       `bson:"-"`
	PurchaseTime  time.Time    `bson:"-"`
	PaymentMethod string       `bson:"-"`
	PaymentID     *string      `bson:"-"`
	CreatedAt     time.Time    `bson:"-"`
	UpdatedAt     time.Time    `bson:"-"`
}

type TicketFilter struct {
	ID            *string
	IDs           []string
	SessionID     *string
	MovieID       *string
	UserID        *string
	SeatNumber    *string
	Status        *TicketStatus
	PaymentMethod *string
}

type TicketUpdateData struct {
	Status        *TicketStatus
	PaymentMethod *string
	PurchaseTime  *time.Time
	PaymentID     *string
	Price         *float64
}

// Helpers for creating pointers
func (ts TicketStatus) Ptr() *TicketStatus { return &ts }
func TimePtr(t time.Time) *time.Time       { return &t }
