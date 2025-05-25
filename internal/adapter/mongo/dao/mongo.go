package dao

import (
	"ap2final_ticket_service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ticket struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SessionID     primitive.ObjectID `bson:"session_id"`
	MovieID       primitive.ObjectID `bson:"movie_id"`
	SeatNumber    string             `bson:"seat_number"`
	Price         float64            `bson:"price"`
	Status        string             `bson:"status"`
	UserID        primitive.ObjectID `bson:"user_id"`
	PurchaseTime  time.Time          `bson:"purchase_time"`
	PaymentMethod string             `bson:"payment_method"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func FromModel(ticket models.Ticket) (Ticket, error) {
	var objID primitive.ObjectID
	var err error

	if ticket.ID != "" {
		objID, err = primitive.ObjectIDFromHex(ticket.ID)
		if err != nil {
			return Ticket{}, err
		}
	}

	sessionID, err := primitive.ObjectIDFromHex(ticket.SessionID)
	if err != nil {
		return Ticket{}, err
	}

	movieID, err := primitive.ObjectIDFromHex(ticket.MovieID)
	if err != nil {
		return Ticket{}, err
	}

	userID, err := primitive.ObjectIDFromHex(ticket.UserID)
	if err != nil {
		return Ticket{}, err
	}

	return Ticket{
		ID:            objID,
		SessionID:     sessionID,
		MovieID:       movieID,
		SeatNumber:    ticket.SeatNumber,
		Price:         ticket.Price,
		Status:        string(ticket.Status),
		UserID:        userID,
		PurchaseTime:  ticket.PurchaseTime,
		PaymentMethod: ticket.PaymentMethod,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}, nil
}

func ToModel(ticket Ticket) models.Ticket {
	return models.Ticket{
		ID:            ticket.ID.Hex(),
		SessionID:     ticket.SessionID.Hex(),
		MovieID:       ticket.MovieID.Hex(),
		SeatNumber:    ticket.SeatNumber,
		Price:         ticket.Price,
		Status:        models.TicketStatus(ticket.Status),
		UserID:        ticket.UserID.Hex(),
		PurchaseTime:  ticket.PurchaseTime,
		PaymentMethod: ticket.PaymentMethod,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}
}

func FromTicketFilter(filter models.TicketFilter) (bson.M, error) {
	query := bson.M{}

	if filter.ID != nil {
		objID, err := primitive.ObjectIDFromHex(*filter.ID)
		if err != nil {
			return query, err
		}
		query["_id"] = objID
	}

	if filter.IDs != nil {
		var objIDs []primitive.ObjectID
		for _, id := range filter.IDs {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return query, err
			}
			objIDs = append(objIDs, objID)
		}
		query["_id"] = bson.M{"$in": objIDs}
	}

	if filter.SessionID != nil {
		sessionID, err := primitive.ObjectIDFromHex(*filter.SessionID)
		if err != nil {
			return query, err
		}
		query["session_id"] = sessionID
	}

	if filter.MovieID != nil {
		movieID, err := primitive.ObjectIDFromHex(*filter.MovieID)
		if err != nil {
			return query, err
		}
		query["movie_id"] = movieID
	}

	if filter.UserID != nil {
		userID, err := primitive.ObjectIDFromHex(*filter.UserID)
		if err != nil {
			return query, err
		}
		query["user_id"] = userID
	}

	if filter.SeatNumber != nil {
		query["seat_number"] = *filter.SeatNumber
	}

	if filter.Status != nil {
		query["status"] = *filter.Status
	}

	if filter.PaymentMethod != nil {
		query["payment_method"] = *filter.PaymentMethod
	}

	return query, nil
}

func FromTicketUpdateData(update models.TicketUpdateData) bson.M {
	query := bson.M{}

	if update.Status != nil {
		query["status"] = *update.Status
	}

	if update.PaymentMethod != nil {
		query["payment_method"] = *update.PaymentMethod
	}

	if update.PurchaseTime != nil {
		query["purchase_time"] = *update.PurchaseTime
	}

	if update.Price != nil {
		query["price"] = *update.Price
	}

	query["updated_at"] = time.Now()

	return bson.M{"$set": query}
}
