package pool

type PoolConfig struct {
	maxRunners int
}

type Option func(*PoolConfig)

func WithMaxRunners(maxRunners int) Option {
	return func(c *PoolConfig) {
		c.maxRunners = maxRunners
	}
}
