package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	caws "gitlab.com/clodoaldomarques/accounting-scripts/pkg/aws"
)

func configure() aws.Config {
	cfg, err := caws.NewCustomConfig(context.TODO())
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	return cfg
}
