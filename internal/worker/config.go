package worker

import (
	"time"
)

type Config struct {
	BatchSize int           `env:"BATCH_SIZE"` // default: 100
	Interval  time.Duration `env:"INTERVAL"`   // default: 1s
	Timeout   time.Duration `env:"TIMEOUT"`    // default: 10s
}

func (c Config) batchSize() int {
	s := c.BatchSize
	if s == 0 {
		s = 100
	}
	return s
}

func (c Config) interval() time.Duration {
	i := c.Interval
	if i == 0 {
		i = time.Second
	}
	return i
}

func (c Config) timeout() time.Duration {
	t := c.Timeout
	if t == 0 {
		t = 10 * time.Second
	}
	return t
}
