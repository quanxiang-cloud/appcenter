package define

import "context"

type Msg struct {
	ctx context.Context

	AppID    string `json:"appID"`
	CreateBy string `json:"createBy"`
	UserName string `json:"userName"`
	Content  int    `json:"content"` // bits of server
}

type Response struct{}
