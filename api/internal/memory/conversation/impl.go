package conversation

import (
	"context"
	"time"

	"github.com/tanmaij/zylo/internal/model"
	"github.com/tanmaij/zylo/pkg/redis"
)

func New() Impl {
	return Impl{}
}

type Impl struct {
	Client *redis.Client
}

func (impl Impl) SetConversation(ctx context.Context, key string, duration time.Duration, value model.Conversation) error {
	return impl.Client.Set(ctx, key, duration, value)
}

func (impl Impl) GetConversation(ctx context.Context, key string) (model.Conversation, error) {
	var rs model.Conversation
	if err := impl.Client.Get(ctx, key, &rs); err != nil {
		return model.Conversation{}, err
	}

	return rs, nil
}
