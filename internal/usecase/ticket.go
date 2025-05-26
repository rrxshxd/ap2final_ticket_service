package usecase

import (
	"context"
	"log/slog"
	"time"

	"ap2final_ticket_service/internal/adapter/cache"
	"ap2final_ticket_service/internal/models"
)

type ticketUseCase struct {
	repo  TicketRepository
	cache cache.TicketCache
	log   *slog.Logger
}

func NewTicketUseCase(repo TicketRepository, cache cache.TicketCache, log *slog.Logger) TicketUseCase {
	return &ticketUseCase{
		repo:  repo,
		cache: cache,
		log:   log,
	}
}

func (uc *ticketUseCase) ReserveTicket(
	ctx context.Context,
	sessionID, movieID, userID, seatNumber string,
	price float64,
) (*models.Ticket, error) {
	if available, err := uc.cache.GetSeatAvailability(ctx, sessionID, seatNumber); err == nil && available != nil {
		if !*available {
			return nil, models.ErrSeatAlreadyTaken
		}
	} else {
		available, err := uc.repo.IsSeatAvailable(ctx, sessionID, seatNumber)
		if err != nil {
			return nil, err
		}
		if !available {
			_ = uc.cache.CacheSeatAvailability(ctx, sessionID, seatNumber, false)
			return nil, models.ErrSeatAlreadyTaken
		}
		_ = uc.cache.CacheSeatAvailability(ctx, sessionID, seatNumber, true)
	}

	ticket := &models.Ticket{
		SessionID:    sessionID,
		MovieID:      movieID,
		UserID:       userID,
		SeatNumber:   seatNumber,
		Price:        price,
		Status:       models.TicketStatusReserved,
		PurchaseTime: time.Time{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdTicket, err := uc.repo.InsertOne(ctx, ticket)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.CacheTicket(ctx, &createdTicket); err != nil {
		uc.log.Warn("failed to cache ticket", "ticket_id", createdTicket.ID, "error", err)
	}

	_ = uc.cache.CacheSeatAvailability(ctx, sessionID, seatNumber, false)

	_ = uc.cache.InvalidateUserTickets(ctx, userID)

	return &createdTicket, nil
}

func (uc *ticketUseCase) ConfirmPayment(
	ctx context.Context,
	ticketID, paymentMethod string,
) (*models.Ticket, error) {
	existing, err := uc.cache.GetTicket(ctx, ticketID)
	if err != nil {
		uc.log.Warn("failed to get ticket from cache", "ticket_id", ticketID, "error", err)
	}

	if existing == nil {
		existingFromDB, err := uc.repo.FindOne(ctx, models.TicketFilter{ID: &ticketID})
		if err != nil {
			return nil, err
		}
		existing = &existingFromDB
	}

	if existing.Status != models.TicketStatusReserved {
		return nil, models.ErrTicketNotReserved
	}

	paymentID := "pay_" + time.Now().Format("20060102150405")

	update := models.TicketUpdateData{
		Status:        models.TicketStatusPaid.Ptr(),
		PaymentMethod: &paymentMethod,
		PaymentID:     &paymentID,
		PurchaseTime:  models.TimePtr(time.Now()),
	}

	updatedTicket, err := uc.repo.UpdateOne(
		ctx,
		models.TicketFilter{ID: &ticketID},
		update,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.CacheTicket(ctx, &updatedTicket); err != nil {
		uc.log.Warn("failed to cache paid ticket", "ticket_id", updatedTicket.ID, "error", err)
	}

	_ = uc.cache.InvalidateUserTickets(ctx, updatedTicket.UserID)

	return &updatedTicket, nil
}

func (uc *ticketUseCase) CancelTicket(ctx context.Context, ticketID string) error {
	existing, err := uc.cache.GetTicket(ctx, ticketID)
	if err != nil {
		uc.log.Warn("failed to get ticket from cache for cancellation", "ticket_id", ticketID, "error", err)
	}

	if existing == nil {
		existingFromDB, err := uc.repo.FindOne(ctx, models.TicketFilter{ID: &ticketID})
		if err != nil {
			return err
		}
		existing = &existingFromDB
	}

	if existing.Status != models.TicketStatusReserved {
		return models.ErrTicketNotReserved
	}

	_, err = uc.repo.UpdateOne(
		ctx,
		models.TicketFilter{ID: &ticketID},
		models.TicketUpdateData{
			Status: models.TicketStatusCancelled.Ptr(),
		},
	)
	if err != nil {
		return err
	}

	_ = uc.cache.InvalidateTicket(ctx, ticketID)

	_ = uc.cache.InvalidateUserTickets(ctx, existing.UserID)

	_ = uc.cache.CacheSeatAvailability(ctx, existing.SessionID, existing.SeatNumber, true)

	return nil
}

func (uc *ticketUseCase) GetTicket(ctx context.Context, id string) (*models.Ticket, error) {
	ticket, err := uc.cache.GetTicket(ctx, id)
	if err != nil {
		uc.log.Warn("failed to get ticket from cache", "ticket_id", id, "error", err)
	}

	if ticket != nil {
		return ticket, nil
	}

	ticketFromDB, err := uc.repo.FindOne(ctx, models.TicketFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	if err := uc.cache.CacheTicket(ctx, &ticketFromDB); err != nil {
		uc.log.Warn("failed to cache ticket after DB fetch", "ticket_id", id, "error", err)
	}

	return &ticketFromDB, nil
}

func (uc *ticketUseCase) GetUserTickets(ctx context.Context, userID string) ([]*models.Ticket, error) {
	tickets, err := uc.cache.GetUserTickets(ctx, userID)
	if err != nil {
		uc.log.Warn("failed to get user tickets from cache", "user_id", userID, "error", err)
	}

	if tickets != nil {
		return tickets, nil
	}

	ticketsFromDB, err := uc.repo.Find(ctx, models.TicketFilter{UserID: &userID})
	if err != nil {
		return nil, err
	}

	result := make([]*models.Ticket, len(ticketsFromDB))
	for i := range ticketsFromDB {
		result[i] = &ticketsFromDB[i]
	}

	if err := uc.cache.CacheUserTickets(ctx, userID, result); err != nil {
		uc.log.Warn("failed to cache user tickets after DB fetch", "user_id", userID, "error", err)
	}

	return result, nil
}

func (uc *ticketUseCase) GetAllTickets(ctx context.Context) ([]*models.Ticket, error) {
	ticketsFromDB, err := uc.repo.Find(ctx, models.TicketFilter{})
	if err != nil {
		return nil, err
	}

	result := make([]*models.Ticket, len(ticketsFromDB))
	for i := range ticketsFromDB {
		result[i] = &ticketsFromDB[i]
	}

	return result, nil
}

func (uc *ticketUseCase) GetMovieTickets(ctx context.Context, movieID string) ([]*models.Ticket, error) {
	ticketsFromDB, err := uc.repo.Find(ctx, models.TicketFilter{MovieID: &movieID})
	if err != nil {
		return nil, err
	}

	result := make([]*models.Ticket, len(ticketsFromDB))
	for i := range ticketsFromDB {
		result[i] = &ticketsFromDB[i]
	}

	return result, nil
}

func (uc *ticketUseCase) CheckSeatAvailability(ctx context.Context, sessionID, seatNumber string) (bool, error) {
	if available, err := uc.cache.GetSeatAvailability(ctx, sessionID, seatNumber); err == nil && available != nil {
		return *available, nil
	}

	available, err := uc.repo.IsSeatAvailable(ctx, sessionID, seatNumber)
	if err != nil {
		return false, err
	}

	_ = uc.cache.CacheSeatAvailability(ctx, sessionID, seatNumber, available)

	return available, nil
}
