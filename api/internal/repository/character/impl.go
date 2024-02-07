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
	return []model.Character{
		{
			Name:              "Dạy làm giàu",
			AvatarURL:         "https://variety.com/wp-content/uploads/2022/11/Elon-Musk-Twitter-CEO.png?w=710",
			Description:       `"Bạn sẽ làm được"`,
			Address:           "Mỹ Tho",
			SystemMsg:         "You only give riddles",
			DefaultHistoryMsg: []model.Message{},
		},
	}, nil
}
