package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	caws "github.com/clodoaldomarques/core-sdk/pkg/aws"
	"github.com/clodoaldomarques/ledger-config/configs"
)

func configure() aws.Config {
	c := configs.New()
	cfg, err := caws.NewCustomConfig(context.TODO(), c.AwsRegion, c.AwsAddress, c.AccessKeyID, c.SecretAccessKey)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	return cfg
}
