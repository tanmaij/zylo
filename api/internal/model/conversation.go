package model

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	AccessToken string `json:"access_token"`

	Character Character `json:"character"`
	Messages  []Message `json:"messages"`
}
