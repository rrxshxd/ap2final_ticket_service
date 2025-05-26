package models

import "errors"

var (
	// Ticket related errors
	ErrTicketNotFound      = errors.New("Ticket not found")
	ErrTicketAlreadyPaid   = errors.New("Ticket already paid")
	ErrTicketAlreadyExists = errors.New("Ticket already exists")
	ErrTicketNotReserved   = errors.New("Ticket not reserved")
	ErrTicketExpired       = errors.New("Ticket reservation has expired")
	ErrTicketCancelled     = errors.New("Ticket has been cancelled")

	// Seat/session related errors
	ErrSeatAlreadyTaken  = errors.New("Seat already taken")
	ErrInvalidSeatNumber = errors.New("Invalid seat number")
	ErrSessionNotFound   = errors.New("Movie session not found")
	ErrSessionNotActive  = errors.New("Movie session not active")
	ErrSessionFull       = errors.New("Movie session is full")

	// Payment related errors
	ErrPaymentFailed        = errors.New("Payment processing failed")
	ErrInsufficientFunds    = errors.New("Insufficient funds for payment")
	ErrInvalidPaymentMethod = errors.New("Invalid payment method")

	// User related errors
	ErrUserNotFound   = errors.New("User not found")
	ErrUserNotAllowed = errors.New("User is not allowed to perform this action")

	// Validation errors
	ErrInvalidTicketData = errors.New("Invalid ticket data")

	// User related errors
	ErrInvalidUserID = errors.New("Invalid user ID")
)
