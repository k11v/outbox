package kafkautil

// Config holds Kafka configuration.
type Config struct {
	Brokers []string `env:"BROKERS,required" envSeparator:","` // required
}
