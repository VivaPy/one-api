package history

import "github.com/songquanpeng/one-api/relay/model"

type MessageToSend struct {
	UserID   int             `json:"userID"`
	Messages []model.Message `json:"messages"`
	Usage    Usage           `json:"usage"`
}

type Usage struct {
	Latency float64
	model.Usage
}
