package worker

import (
	"time"
)

type Config struct {
	BatchSize int           `env:"BATCH_SIZE"` // default: 100
	Interval  time.Duration `env:"INTERVAL"`   // default: 1s
}

func (c Config) BatchSizeReal() int {
	s := c.BatchSize
	if s == 0 {
		s = 100
	}
	return s
}

func (c Config) IntervalReal() time.Duration {
	i := c.Interval
	if i == 0 {
		i = time.Second
	}
	return i
}
