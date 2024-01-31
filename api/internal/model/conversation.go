package model

type Message struct {
	Role string `json:"role"`
	Msg  string `json:"msg"`
}

type Conversation struct {
	AccessToken string `json:"access_token"`

	Character Character `json:"character"`
	Messages  []Message `json:"messages"`
}
