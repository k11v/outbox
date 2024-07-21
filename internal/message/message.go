package message

type Message struct {
	Key   []byte
	Value []byte
}

type Producer struct{}

func (p *Producer) Produce(messages []Message) error {
	return nil
}
