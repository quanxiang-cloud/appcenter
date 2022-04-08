package define

// Msg msg
type Msg struct {
	AppID    string `json:"appID"`
	CreateBy string `json:"createBy"`
	Content  int    `json:"content"` // bits of server
	Ret      int    `json:"-"`
}

// Response response
type Response struct{}
