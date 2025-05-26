package grpc

import (
	"ap2final_ticket_service/internal/adapter/grpc/dto"
	"context"
	"github.com/sorawaslocked/ap2final_protos_gen/base"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/ticket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type TicketServer struct {
	uc  TicketUseCase
	log *slog.Logger
	svc.UnimplementedTicketServiceServer
}

func NewTicketServer(
	uc TicketUseCase,
	log *slog.Logger,
) *TicketServer {
	return &TicketServer{
		uc:  uc,
		log: log,
	}
}

func (s *TicketServer) Create(ctx context.Context, req *svc.CreateRequest) (*svc.CreateResponse, error) {
	createdTicket, err := s.uc.ReserveTicket(
		ctx,
		req.ShowtimeID,
		req.MovieID,
		req.UserID,
		req.SeatNumber,
		req.Price,
	)
	if err != nil {
		s.logError("create", err)
		return nil, dto.FromError(err)
	}

	return &svc.CreateResponse{
		Ticket: dto.FromTicketToPb(*createdTicket),
	}, nil
}

func (s *TicketServer) Get(ctx context.Context, req *svc.GetRequest) (*svc.GetResponse, error) {
	ticket, err := s.uc.GetTicket(ctx, req.ID)
	if err != nil {
		s.logError("get", err)
		return nil, dto.FromError(err)
	}

	return &svc.GetResponse{
		Ticket: dto.FromTicketToPb(*ticket),
	}, nil
}

func (s *TicketServer) GetAll(ctx context.Context, req *svc.GetAllRequest) (*svc.GetAllResponse, error) {
	tickets, err := s.uc.GetAllTickets(ctx)
	if err != nil {
		s.logError("get all", err)
		return nil, dto.FromError(err)
	}

	var ticketsPb []*base.Ticket
	for _, ticket := range tickets {
		ticketsPb = append(ticketsPb, dto.FromTicketToPb(*ticket))
	}

	return &svc.GetAllResponse{
		Tickets: ticketsPb,
	}, nil
}

func (s *TicketServer) GetByUser(ctx context.Context, req *svc.GetByUserRequest) (*svc.GetByUserResponse, error) {
	tickets, err := s.uc.GetUserTickets(ctx, req.UserID)
	if err != nil {
		s.logError("get by user", err)
		return nil, dto.FromError(err)
	}

	var ticketsPb []*base.Ticket
	for _, ticket := range tickets {
		ticketsPb = append(ticketsPb, dto.FromTicketToPb(*ticket))
	}

	return &svc.GetByUserResponse{
		Tickets: ticketsPb,
	}, nil
}

func (s *TicketServer) GetByMovie(ctx context.Context, req *svc.GetByMovieRequest) (*svc.GetByMovieResponse, error) {
	tickets, err := s.uc.GetMovieTickets(ctx, req.MovieID)
	if err != nil {
		s.logError("get by movie", err)
		return nil, dto.FromError(err)
	}

	var ticketsPb []*base.Ticket
	for _, ticket := range tickets {
		ticketsPb = append(ticketsPb, dto.FromTicketToPb(*ticket))
	}

	return &svc.GetByMovieResponse{
		Tickets: ticketsPb,
	}, nil
}

func (s *TicketServer) Update(ctx context.Context, req *svc.UpdateRequest) (*svc.UpdateResponse, error) {
	if req.Status != nil && *req.Status == "PAID" {
		updatedTicket, err := s.uc.ConfirmPayment(ctx, req.ID, "card")
		if err != nil {
			s.logError("update", err)
			return nil, dto.FromError(err)
		}

		return &svc.UpdateResponse{
			Ticket: dto.FromTicketToPb(*updatedTicket),
		}, nil
	}

	if req.Status != nil && *req.Status == "CANCELLED" {
		err := s.uc.CancelTicket(ctx, req.ID)
		if err != nil {
			s.logError("update", err)
			return nil, dto.FromError(err)
		}

		ticket, err := s.uc.GetTicket(ctx, req.ID)
		if err != nil {
			s.logError("update", err)
			return nil, dto.FromError(err)
		}

		return &svc.UpdateResponse{
			Ticket: dto.FromTicketToPb(*ticket),
		}, nil
	}

	s.log.Warn("Update method not fully implemented for this operation")
	return nil, status.Error(codes.Unimplemented, "Update method not fully implemented")
}

func (s *TicketServer) Delete(ctx context.Context, req *svc.DeleteRequest) (*svc.DeleteResponse, error) {
	ticket, err := s.uc.GetTicket(ctx, req.ID)
	if err != nil {
		s.logError("delete", err)
		return nil, dto.FromError(err)
	}

	err = s.uc.CancelTicket(ctx, req.ID)
	if err != nil {
		s.logError("delete", err)
		return nil, dto.FromError(err)
	}

	return &svc.DeleteResponse{
		Ticket: dto.FromTicketToPb(*ticket),
	}, nil
}
