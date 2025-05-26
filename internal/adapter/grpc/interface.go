package grpc

import (
	"ap2final_ticket_service/internal/models"
	"context"
)

type TicketUseCase interface {
	ReserveTicket(ctx context.Context, sessionID, movieID, userID, seatNumber string, price float64) (*models.Ticket, error)
	ConfirmPayment(ctx context.Context, ticketID, paymentMethod string) (*models.Ticket, error)
	CancelTicket(ctx context.Context, ticketID string) error
	GetTicket(ctx context.Context, id string) (*models.Ticket, error)
	GetUserTickets(ctx context.Context, userID string) ([]*models.Ticket, error)
	GetAllTickets(ctx context.Context) ([]*models.Ticket, error)
	GetMovieTickets(ctx context.Context, movieID string) ([]*models.Ticket, error)
	CheckSeatAvailability(ctx context.Context, sessionID, seatNumber string) (bool, error)
}
