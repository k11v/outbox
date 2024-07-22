package message

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides an interface for storing and retrieving messages.
type Store interface {
	// Add adds one or more messages to the store.
	Add(ctx context.Context, messages ...Message) error

	// GetUndelivered returns a list of undelivered messages.
	GetUndelivered(ctx context.Context) ([]Message, error)

	// SetDelivered marks the messages as delivered.
	SetDelivered(ctx context.Context, ids ...uuid.UUID) error
}

// PostgresStore implements the Store interface using a PostgreSQL database.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

func (s *PostgresStore) Add(ctx context.Context, messages ...Message) error {
	return nil
}

func (s *PostgresStore) GetUndelivered(ctx context.Context) ([]Message, error) {
	return nil, nil
}

func (s *PostgresStore) SetDelivered(ctx context.Context, ids ...uuid.UUID) error {
	return nil
}
