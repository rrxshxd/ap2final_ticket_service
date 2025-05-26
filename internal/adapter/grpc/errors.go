package grpc

import (
	"ap2final_ticket_service/internal/models"
	"errors"
	"fmt"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
)

func (s *TicketServer) logError(op string, err error) {
	if !errors.Is(err, models.ErrTicketNotFound) {
		s.log.Error(fmt.Sprintf("ticket %s", op), logger.Err(err))
	}
}
