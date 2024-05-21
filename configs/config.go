package configs

import "sync"

const (
	ConfigKey = "config"
)

var (
	configOnce sync.Once
	inst       *Config
)

type Config struct {
	GitExec string
}

type Option func(*Config)

func WithGitExec(gitExec string) Option {
	return func(c *Config) {
		c.GitExec = gitExec
	}
}

func NewConfig(options ...Option) {
	configOnce.Do(func() {
		config := &Config{}
		for _, option := range options {
			option(config)
		}
		inst = config
	})
}

func Instance() *Config {
	return inst
}
