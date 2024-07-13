package rpc

import (
	"context"

	"github.com/klaidliadon/trace/proto"

	"github.com/go-chi/httplog/v2"
)

var _ proto.Main = Main{}

type Main struct {
	Secondary proto.Secondary
}

func (s Main) Do(ctx context.Context) (bool, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("main implementation...")
	if _, err := s.Secondary.Do(ctx); err != nil {
		return false, err
	}
	logger.Info("main implementation 2")
	return true, nil
}

var _ proto.Secondary = Secondary{}

type Secondary struct {
	Tertiary proto.Tertiary
}

func (s Secondary) Do(ctx context.Context) (bool, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("secondary implementation...")
	if _, err := s.Tertiary.Do(ctx); err != nil {
		return false, err
	}
	logger.Info("secondary implementation complete!")
	return true, nil
}

var _ proto.Tertiary = Tertiary{}

type Tertiary struct{}

func (s Tertiary) Do(ctx context.Context) (bool, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("tertiary implementation...")
	logger.Info("tertiary implementation 2")
	return true, nil
}
