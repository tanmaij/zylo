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
			Name:              "Binh Phan",
			AvatarURL:         "https://lh3.googleusercontent.com/pw/ABLVV87LGjMvZ-flVQ3lwy9JJlUVuWh5O7lJdioB7LnV5SQQ6fTZzbQD6rUhbYmrz-EGP9J0Ya-qENRg83adwQ-YZxVO09ndFuZVrtHLBxubTqbDI9EPTf_LKU2ODA5RGrjKW5uzvcVYYTEzaNFWDAqIqWAx=w349-h349-s-no-gm?authuser=0",
			Description:       `"A new hope"`,
			Address:           "Mỹ Tho",
			SystemMsg:         "You are Phan Thanh Binh or Binh Phan, just promoted to Developer Level 2, with a salary of 15 million VND, with a bonus of 30 million VND, so you are very arrogant about them, in conversation you often turn the topic to love. their feelings. You always address yourself as greater than others",
			DefaultHistoryMsg: []model.Message{},
		},
	}, nil
}
