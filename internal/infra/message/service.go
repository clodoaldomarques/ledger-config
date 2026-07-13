package message

import (
	"context"
	"time"

	"github.com/clodoaldomarques/core-sdk/pkg/logger"
	"github.com/clodoaldomarques/core-sdk/pkg/sns"
	"github.com/clodoaldomarques/ledger-config/config"
	"github.com/clodoaldomarques/ledger-config/internal/domain/ledger"
	"github.com/google/uuid"
)

type Topic struct {
	p *sns.Publisher
}

func New(ctx context.Context) *Topic {
	return &Topic{
		p: sns.NewPublisher(ctx, config.New()),
	}
}

func (t Topic) Emit(ctx context.Context, cid string, c ledger.Config) error {
	evt := sns.Event{
		EventID:   uuid.New(),
		EventType: "ledger",
		EventData: ToConfigMessage(c),
		EventDate: time.Now(),
	}
	return t.p.Emit(ctx, evt)
}

func (t Topic) Close(ctx context.Context) {
	logger.Info(ctx, "ending topic connection", logger.Fields{})
}
