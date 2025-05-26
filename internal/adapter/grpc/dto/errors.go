package dto

import (
	"ap2final_ticket_service/internal/models"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromError(err error) error {
	if errors.Is(err, models.ErrTicketNotFound) {
		return status.Error(codes.NotFound, "ticket not found")
	}

	//if errors.Is(err, models.ErrTicketAlreadyExists) {
	//	return status.Error(codes.AlreadyExists, "ticket already exists")
	//}

	if errors.Is(err, models.ErrInvalidTicketData) {
		return status.Error(codes.InvalidArgument, "invalid input")
	}

	return status.Error(codes.Internal, "internal server error")
}
