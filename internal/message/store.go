package message

import (
	"github.com/google/uuid"
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
