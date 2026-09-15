package message

import (
	"context"
	"time"

	"github.com/clodoaldomarques/core-sdk/pkg/aws/sns"
	"github.com/clodoaldomarques/core-sdk/pkg/otel/logger"
	"github.com/clodoaldomarques/core-sdk/pkg/otel/tracer"
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
	span, ctx := tracer.NewSpanFromContext(ctx, "Topic::Emit", map[string]any{
		"cid":    cid,
		"config": c,
	})
	defer span.End()

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
