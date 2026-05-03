package configs

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	AppPort         int
	AwsAddress      string
	AwsRegion       string
	AccessKeyID     string
	SecretAccessKey string
	GriphookIAUrl   string
}

type Option func(*Config)

var (
	singleton sync.Once
	instance  *Config
)

func New(options ...Option) *Config {
	singleton.Do(func() {
		instance = &Config{
			AppPort:         GetInt("APP_PORT", 8080),
			AwsAddress:      GetString("AWS_ADDRESS", "http://192.168.49.2:30002"),
			AwsRegion:       GetString("AWS_REGION", "us-east-1"),
			AccessKeyID:     GetString("AWS_ACCESS_KEY_ID", "test"),
			SecretAccessKey: GetString("AWS_SECRET_ACCESS_KEY", "test"),
			GriphookIAUrl:   GetString("IA_URL", "http://192.168.49.2:30007"),
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
