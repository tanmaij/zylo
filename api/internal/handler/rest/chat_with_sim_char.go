package rest

import chatwithsimchar "github.com/tanmaij/zylo/internal/controller/chat_with_sim_char"

type Impl struct {
	ctrl chatwithsimchar.Impl
}

func New(ctrl chatwithsimchar.Impl) Impl {
	return Impl{ctrl: ctrl}
}
