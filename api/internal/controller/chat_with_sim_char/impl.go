package chatwithsimchar

import (
	conversationRedisRepo "github.com/tanmaij/zylo/internal/memory/conversation"
	characterRepo "github.com/tanmaij/zylo/internal/repository/character"
	"github.com/tanmaij/zylo/pkg/ws"
)

type Impl struct {
	Broadcast ws.Broadcast

	characterRepo         characterRepo.Impl
	conversationRedisRepo conversationRedisRepo.Impl
}

func New(wsBroadcast ws.Broadcast, characterRepo characterRepo.Impl, conversationRedisRepo conversationRedisRepo.Impl) Impl {
	return Impl{wsBroadcast, characterRepo, conversationRedisRepo}
}

type TestInput struct {
	Message string
}

func (i Impl) Test(client ws.Client, inp TestInput) {
	i.Broadcast.EmitToClientUUID(client.UUID, ws.Message{
		EventName: "successfully",
		Data:      "You're the best",
	})
}
