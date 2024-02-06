package conversation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tanmaij/zylo/internal/model"
	"github.com/tanmaij/zylo/pkg/redis"
)

func New(redisClient redis.Client) Impl {
	return Impl{redisClient: redisClient}
}

type Impl struct {
	redisClient redis.Client
}

func genConversationKey(uuid string) string {
	return fmt.Sprintf("%s:%s", "client", uuid)
}

func (impl Impl) Get(ctx context.Context, uuid string) (model.Conversation, error) {
	var rs model.Conversation

	if err := impl.redisClient.Get(ctx, genConversationKey(uuid), &rs); err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return model.Conversation{}, ErrNotFound
		}
		return model.Conversation{}, err
	}

	return rs, nil
}

func (impl Impl) Set(ctx context.Context, uuid string, duration time.Duration, data model.Conversation) error {
	return impl.redisClient.Set(ctx, genConversationKey(uuid), duration, data)
}
