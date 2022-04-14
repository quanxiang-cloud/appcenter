package exec

import (
	"context"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/handle"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	reload = "init-reload"
)

type listReq struct {
	Status int `json:"status"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
}

type listResp struct {
	Total int    `json:"total_count"`
	Data  []*app `json:"data"`
}

type app struct {
	ID       string `json:"id"`
	CreateBy string `json:"createBy"`
}

// define
var (
	requestID interface{} = "Request-Id"
	initChaos interface{} = "init-chaos"
)

// InitExec InitExec
func InitExec(handler *handle.TaskHandler) error {
	c := handler.Config
	cli := client.New(c.InternalNet)
	ctx := context.WithValue(context.Background(), requestID, initChaos)

	req := &listReq{
		Status: -5,
		Page:   0,
		Limit:  100,
	}
	resp := &listResp{}
	for resp.Total > req.Limit*req.Page || req.Page == 0 {
		req.Page++
		if err := client.POST(ctx, &cli, c.KV[reload], req, resp); err != nil {
			return err
		}

		for _, app := range resp.Data {
			handler.Put(ctx, define.Msg{
				AppID:    app.ID,
				CreateBy: app.CreateBy,
			})
		}

		resp.Data = resp.Data[:0]
	}

	return nil
}
