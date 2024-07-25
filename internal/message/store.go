package message

import (
	"context"
	"github.com/k11v/outbox/internal/outbox"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides an interface for storing and retrieving messages.
type Store interface {
	// AddWithOutbox adds messages to the store and outbox.
	//
	// Outbox messages are separate from business messages to allow for design
	// pattern extensibility. For example, you might want to add a message to
	// the outbox when a user signs up.
	AddWithOutbox(ctx context.Context, messages []Message, outboxMessages []outbox.Message) error
}

// PostgresStore implements the Store interface using a PostgreSQL database.
type PostgresStore struct {
	Pool *pgxpool.Pool // required
}
