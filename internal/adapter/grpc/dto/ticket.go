package dto

import (
	"ap2final_ticket_service/internal/models"
	"github.com/sorawaslocked/ap2final_protos_gen/base"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/ticket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToTicketFromCreateRequest(req *svc.CreateRequest) models.Ticket {
	return models.Ticket{
		UserID:        req.UserID,
		MovieID:       req.MovieID,
		SessionID:     req.ShowtimeID,
		SeatNumber:    req.SeatNumber,
		Price:         req.Price,
		Status:        models.TicketStatus(req.Status),
		PaymentMethod: "",
	}
}

func ToTicketUpdateFromUpdateRequest(req *svc.UpdateRequest) (string, models.TicketUpdateData) {
	updateData := models.TicketUpdateData{}

	if req.Status != nil {
		status := models.TicketStatus(*req.Status)
		updateData.Status = &status
	}

	if req.Price != nil {
		updateData.Price = req.Price
	}

	return req.ID, updateData
}

func FromTicketToPb(ticket models.Ticket) *base.Ticket {
	return &base.Ticket{
		ID:         ticket.ID,
		UserID:     ticket.UserID,
		MovieID:    ticket.MovieID,
		ShowtimeID: ticket.SessionID,
		SeatNumber: ticket.SeatNumber,
		Price:      ticket.Price,
		Status:     string(ticket.Status),
		CreatedAt:  timestamppb.New(ticket.CreatedAt),
		UpdatedAt:  timestamppb.New(ticket.UpdatedAt),
		IsDeleted:  false,
	}
}
