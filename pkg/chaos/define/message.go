package define

// Msg msg
type Msg struct {
	AppID    string `json:"appID"`
	CreateBy string `json:"createBy"`
	Content  int    `json:"content"` // bits of server
}

// Response response
type Response struct{}
