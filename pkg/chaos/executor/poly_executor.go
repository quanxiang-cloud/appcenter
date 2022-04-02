package exec

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

// Key
const (
	PolyURL = "polyurl"
)

const (
	perInitTypes = 1
	name         = "全部权限"
	description  = "系统默认角色"
)

// PolyExecutor PolyExecutor
type PolyExecutor struct {
	Client  http.Client
	PolyURL string
}

// PolyReq PolyReq
type PolyReq struct {
	AppID       string      `json:"appID"`
	Scopes      []*ScopesVO `json:"scopes"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Types       int64       `json:"types"`
}

// ScopesVO ScopesVO
type ScopesVO struct {
	Type int16  `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PolyResp PolyResp
type PolyResp = define.Response

// Exec Exec
func (s *PolyExecutor) Exec(ctx context.Context, m define.Msg) error {
	req := &PolyReq{
		AppID: m.AppID,
		Scopes: []*ScopesVO{
			{
				Type: 1,
				ID:   m.CreateBy,
				Name: m.UserName,
			},
		},
		Name:        name,
		Description: description,
		Types:       perInitTypes,
	}

	resp := &PolyResp{}
	err := client.POST(ctx, &s.Client, s.PolyURL, req, resp)
	if err != nil {
		return err
	}
	return nil
}

// Bit Bit
func (*PolyExecutor) Bit() int {
	return define.BitPolyAPI
}
