package define

type Msg struct {
	AppID    string `json:"appID"`
	CreateBy string `json:"createBy"`
	UserName string `json:"userName"`
	Content  int    `json:"content"` // bits of server
}

type Response struct{}
