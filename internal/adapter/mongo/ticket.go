package mongo

import (
	"ap2final_ticket_service/internal/adapter/mongo/dao"
	"ap2final_ticket_service/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionTickets = "tickets"

type Ticket struct {
	col *mongo.Collection
}

func NewTicket(conn *mongo.Database) *Ticket {
	collection := conn.Collection(collectionTickets)

	return &Ticket{col: collection}
}

func (db *Ticket) InsertOne(ctx context.Context, ticket *models.Ticket) (models.Ticket, error) {
	available, err := db.IsSeatAvailable(ctx, ticket.SessionID, ticket.SeatNumber)
	if err != nil {
		return models.Ticket{}, err
	}
	if !available {
		return models.Ticket{}, models.ErrSeatAlreadyTaken
	}

	ticketDao, err := dao.FromModel(*ticket)
	if err != nil {
		return models.Ticket{}, mongoError("primitive.ObjectIDFromHex", err)
	}

	res, err := db.col.InsertOne(ctx, ticketDao)

	if err != nil {
		return models.Ticket{}, mongoError("insertOne", err)
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()

	return db.FindOne(ctx, models.TicketFilter{ID: &id})
}

func (db *Ticket) FindOne(ctx context.Context, filter models.TicketFilter) (models.Ticket, error) {
	var ticketDao dao.Ticket

	query, err := dao.FromTicketFilter(filter)
	if err != nil {
		return models.Ticket{}, mongoError("primitive.ObjectIDFromHex", err)
	}

	err = db.col.FindOne(ctx, query).Decode(&ticketDao)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Ticket{}, models.ErrTicketNotFound
		}

		return models.Ticket{}, mongoError("FindOne", err)
	}

	return dao.ToModel(ticketDao), nil
}

func (db *Ticket) Find(ctx context.Context, filter models.TicketFilter) ([]models.Ticket, error) {
	var ticketDaos []dao.Ticket
	query, err := dao.FromTicketFilter(filter)
	if err != nil {
		return []models.Ticket{}, mongoError("primitive.ObjectIDFromHex", err)
	}

	cur, err := db.col.Find(ctx, query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []models.Ticket{}, models.ErrTicketNotFound
		}

		return []models.Ticket{}, mongoError("Find", err)
	}

	if err = cur.All(ctx, &ticketDaos); err != nil {
		return []models.Ticket{}, mongoError("Cursor.All", err)
	}

	tickets := make([]models.Ticket, len(ticketDaos))

	for i := range tickets {
		tickets[i] = dao.ToModel(ticketDaos[i])
	}

	return tickets, nil
}

func (db *Ticket) UpdateOne(ctx context.Context, filter models.TicketFilter, update models.TicketUpdateData) (models.Ticket, error) {
	query, err := dao.FromTicketFilter(filter)

	if err != nil {
		return models.Ticket{}, mongoError("primitive.ObjectIDFromHex", err)
	}

	res, err := db.col.UpdateOne(ctx, query, dao.FromTicketUpdateData(update))
	if err != nil {
		return models.Ticket{}, mongoError("UpdateOne", err)
	}

	if res.MatchedCount == 0 {
		return models.Ticket{}, models.ErrTicketNotFound
	}

	return db.FindOne(ctx, filter)
}

func (db *Ticket) DeleteOne(ctx context.Context, filter models.TicketFilter) (models.Ticket, error) {
	var ticketDao dao.Ticket

	query, err := dao.FromTicketFilter(filter)
	if err != nil {
		return models.Ticket{}, mongoError("primitive.ObjectIDFromHex", err)
	}

	err = db.col.FindOne(ctx, query).Decode(&ticketDao)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Ticket{}, models.ErrTicketNotFound
		}

		return models.Ticket{}, mongoError("FindOne", err)
	}

	res, err := db.col.DeleteOne(ctx, query)
	if err != nil {
		return models.Ticket{}, mongoError("DeleteOne", err)
	}

	if res.DeletedCount == 0 {
		return models.Ticket{}, models.ErrTicketNotFound
	}

	return dao.ToModel(ticketDao), err
}

func (db *Ticket) IsSeatAvailable(ctx context.Context, sessionID string, seatNumber string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return false, mongoError("primitive.ObjectIDFromHex", err)
	}

	filter := bson.M{
		"session_id":  objID,
		"seat_number": seatNumber,
		"status": bson.M{
			"$in": []string{
				string(models.TicketStatusReserved),
				string(models.TicketStatusPaid),
			},
		},
	}

	count, err := db.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, mongoError("CountDocuments", err)
	}

	return count == 0, nil
}

func (db *Ticket) InsertMany(ctx context.Context, tickets []models.Ticket) ([]models.Ticket, error) {
	var daoModels []interface{}
	ids := make([]string, 0, len(tickets))

	for _, ticket := range tickets {
		daoModel, err := dao.FromModel(ticket)
		if err != nil {
			return nil, mongoError("FromModel", err)
		}
		daoModels = append(daoModels, daoModel)
		ids = append(ids, ticket.ID)
	}

	_, err := db.col.InsertMany(ctx, daoModels)
	if err != nil {
		return nil, mongoError("InsertMany", err)
	}

	return db.Find(ctx, models.TicketFilter{IDs: ids})
}

func (db *Ticket) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	session, err := db.col.Database().Client().StartSession()
	if err != nil {
		return mongoError("StartSession", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})
	return err
}

func (db *Ticket) UpdateStatus(ctx context.Context, ticketID string, newStatus models.TicketStatus) (models.Ticket, error) {
	update := models.TicketUpdateData{
		Status: &newStatus,
	}
	filter := models.TicketFilter{ID: &ticketID}
	return db.UpdateOne(ctx, filter, update)
}
