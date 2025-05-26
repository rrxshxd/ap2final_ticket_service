package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ap2final_ticket_service/internal/config"
	"ap2final_ticket_service/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(cfg config.Redis) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisCache{
		client: rdb,
		ttl:    cfg.TTL,
	}
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) CacheTicket(ctx context.Context, ticket *models.Ticket) error {
	key := r.ticketKey(ticket.ID)

	data, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %w", err)
	}

	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache ticket: %w", err)
	}

	return nil
}

func (r *RedisCache) GetTicket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	key := r.ticketKey(ticketID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get ticket from cache: %w", err)
	}

	var ticket models.Ticket
	err = json.Unmarshal([]byte(data), &ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket: %w", err)
	}

	return &ticket, nil
}

func (r *RedisCache) CacheUserTickets(ctx context.Context, userID string, tickets []*models.Ticket) error {
	key := r.userTicketsKey(userID)

	data, err := json.Marshal(tickets)
	if err != nil {
		return fmt.Errorf("failed to marshal user tickets: %w", err)
	}

	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache user tickets: %w", err)
	}

	return nil
}

func (r *RedisCache) GetUserTickets(ctx context.Context, userID string) ([]*models.Ticket, error) {
	key := r.userTicketsKey(userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get user tickets from cache: %w", err)
	}

	var tickets []*models.Ticket
	err = json.Unmarshal([]byte(data), &tickets)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user tickets: %w", err)
	}

	return tickets, nil
}

func (r *RedisCache) InvalidateTicket(ctx context.Context, ticketID string) error {
	key := r.ticketKey(ticketID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) InvalidateUserTickets(ctx context.Context, userID string) error {
	key := r.userTicketsKey(userID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) CacheSeatAvailability(ctx context.Context, sessionID, seatNumber string, available bool) error {
	key := r.seatKey(sessionID, seatNumber)

	err := r.client.Set(ctx, key, available, 10*time.Minute).Err() // Shorter TTL for seat availability
	if err != nil {
		return fmt.Errorf("failed to cache seat availability: %w", err)
	}

	return nil
}

func (r *RedisCache) GetSeatAvailability(ctx context.Context, sessionID, seatNumber string) (*bool, error) {
	key := r.seatKey(sessionID, seatNumber)

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get seat availability from cache: %w", err)
	}

	available := result == "true"
	return &available, nil
}

func (r *RedisCache) ticketKey(ticketID string) string {
	return fmt.Sprintf("ticket:%s", ticketID)
}

func (r *RedisCache) userTicketsKey(userID string) string {
	return fmt.Sprintf("user_tickets:%s", userID)
}

func (r *RedisCache) seatKey(sessionID, seatNumber string) string {
	return fmt.Sprintf("seat:%s:%s", sessionID, seatNumber)
}
