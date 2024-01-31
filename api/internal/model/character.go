package model

type Character struct {
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
	Address     string `json:"address"`
	SystemMsg   string `json:"system_msg"`
}
