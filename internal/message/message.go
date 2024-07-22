package message

import "github.com/google/uuid"

type Message struct {
	Topic string
	Key   []byte
	Value []byte

	id uuid.UUID // empty for new messages
}
