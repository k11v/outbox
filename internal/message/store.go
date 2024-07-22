package message

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides an interface for storing and retrieving messages.
type Store interface {
	// Add adds one or more messages to the store.
	Add(messages ...Message) error

	// GetUndelivered returns a list of undelivered messages.
	GetUndelivered() ([]Message, error)

	// SetDelivered marks the messages as delivered.
	SetDelivered(ids ...uuid.UUID) error
}

// PostgresStore implements the Store interface using a PostgreSQL database.
type PostgresStore struct {
	Pool *pgxpool.Pool
}

func (s *PostgresStore) Add(messages ...Message) error {
	return nil
}

func (s *PostgresStore) GetUndelivered() ([]Message, error) {
	return nil, nil
}

func (s *PostgresStore) SetDelivered(ids ...uuid.UUID) error {
	return nil
}
