package cache

import (
	"ap2final_ticket_service/internal/models"
	"context"
)

type TicketCache interface {
	// Ticket caching
	CacheTicket(ctx context.Context, ticket *models.Ticket) error
	GetTicket(ctx context.Context, ticketID string) (*models.Ticket, error)
	InvalidateTicket(ctx context.Context, ticketID string) error

	// User tickets caching
	CacheUserTickets(ctx context.Context, userID string, tickets []*models.Ticket) error
	GetUserTickets(ctx context.Context, userID string) ([]*models.Ticket, error)
	InvalidateUserTickets(ctx context.Context, userID string) error

	// Seat availability caching
	CacheSeatAvailability(ctx context.Context, sessionID, seatNumber string, available bool) error
	GetSeatAvailability(ctx context.Context, sessionID, seatNumber string) (*bool, error)

	// Health check
	Ping(ctx context.Context) error
	Close() error
}
