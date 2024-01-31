package character

import (
	"context"

	"github.com/tanmaij/zylo/internal/model"
)

type Impl struct {
}

func New() Impl {
	return Impl{}
}

func (impl Impl) List(ctx context.Context) ([]model.Character, error) {
	return []model.Character{}, nil
}
