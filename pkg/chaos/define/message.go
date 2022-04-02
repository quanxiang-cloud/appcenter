package define

// Msg msg
type Msg struct {
	AppID    string `json:"appID"`
	CreateBy string `json:"createBy"`
	UserName string `json:"userName"`
	Content  int    `json:"content"` // bits of server
}

// Response response
type Response struct{}
