package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	AppPort            int
	AwsAddress         string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
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
			AppPort:            GetInt("APP_PORT", 8080),
			AwsAddress:         GetString("AWS_ADDRESS", ""),
			AwsRegion:          GetString("AWS_REGION", ""),
			AwsAccessKeyID:     GetString("AWS_ACCESS_KEY_ID", ""),
			AwsSecretAccessKey: GetString("AWS_SECRET_ACCESS_KEY", ""),
			GriphookIAUrl:      GetString("IA_URL", ""),
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

func GetString(env string, def string) string {
	if e := os.Getenv(env); e != "" {
		return e
	}
	return def
}

func GetInt(env string, def int) int {
	i, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return def
	}
	return i
}
