package client

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	ChaosURL = "http://localhost:6666/init"
)

func NewChaos(c *config.Configs) Chaos {
	return &chaos{
		client: client.New(c.InternalNet),
	}
}

type chaos struct {
	client http.Client
}

type Chaos interface {
	Init(ctx context.Context, appID, createBy, userName string, content int) error
}

type Req = define.Msg

type Resp = define.Response

func (c *chaos) Init(ctx context.Context, appID, createBy, userName string, content int) error {
	req := &Req{
		AppID:    appID,
		CreateBy: createBy,
		UserName: userName,
		Content:  content,
	}
	return client.POST(ctx, &c.client, ChaosURL, req, &Resp{})
}
