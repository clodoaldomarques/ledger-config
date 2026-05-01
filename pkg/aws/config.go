package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"gitlab.com/clodoaldomarques/accounting-scripts/configs"
)

func NewCustomCredentials(c *configs.Config) aws.CredentialsProvider {
	return aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		creds := aws.Credentials{
			AccessKeyID:     c.AccessKeyID,
			SecretAccessKey: c.SecretAccessKey,
			Source:          "Environment",
		}
		if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
			return aws.Credentials{}, fmt.Errorf("credenciais AWS ausentes nas variáveis de ambiente")
		}
		return creds, nil
	}))
}

func NewCustomConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(configs.New().AwsRegion),
		config.WithBaseEndpoint(configs.New().AwsAddress),
		config.WithCredentialsProvider(NewCustomCredentials(configs.New())),
	)

	if err != nil {
		return aws.Config{}, fmt.Errorf("falha ao carregar configuração AWS: %w", err)
	}

	return cfg, nil
}
