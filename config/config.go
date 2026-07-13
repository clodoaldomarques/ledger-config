package config

import (
	"sync"

	"github.com/clodoaldomarques/core-sdk/pkg/env"
)

type Config struct {
	AppPort            int
	AwsAddress         string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
	ConfigTopic        string
	GriphookIAUrl      string
}

type Option func(*Config)

var (
	singleton sync.Once
	instance  *Config
)

func New(options ...Option) *Config {
	singleton.Do(func() {
		instance = &Config{
			AppPort:            env.GetInt("APP_PORT", 8080),
			AwsAddress:         env.GetString("AWS_ADDRESS", ""),
			AwsRegion:          env.GetString("AWS_REGION", ""),
			AwsAccessKeyID:     env.GetString("AWS_ACCESS_KEY_ID", ""),
			AwsSecretAccessKey: env.GetString("AWS_SECRET_ACCESS_KEY", ""),
			ConfigTopic:        env.GetString("CONFIG_SNS_TOPIC", ""),
			GriphookIAUrl:      env.GetString("IA_URL", ""),
		}
	})

	for _, optFunc := range options {
		optFunc(instance)
	}

	return instance
}

func WithAppPort(appPort int) Option {
	return func(c *Config) {
		c.AppPort = appPort
	}
}

func WithAwsAddress(AwsAddress string) Option {
	return func(c *Config) {
		c.AwsAddress = AwsAddress
	}
}
func WithAwsRegion(AwsRegion string) Option {
	return func(c *Config) {
		c.AwsRegion = AwsRegion
	}
}

func (c Config) Region() string {
	return c.AwsRegion
}

func (c Config) Address() string {
	return c.AwsAddress
}
func (c Config) AccessKeyID() string {
	return c.AwsAccessKeyID
}
func (c Config) SecretAccessKey() string {
	return c.AwsSecretAccessKey
}

func (c Config) TopicARN() string {
	return c.ConfigTopic
}
