package chaos

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/handle"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// Chaos Chaos
type Chaos struct {
	log     logger.AdaptedLogger
	handler *handle.TaskHandler
}

type InitReq struct {
	Status int `json:"status"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
}

// New New
func New(c *config.Configs, handler *handle.TaskHandler, log logger.AdaptedLogger) *Chaos {
	chaos := &Chaos{
		log:     log,
		handler: handler,
	}
	chaos.handler.Run()

	return chaos
}

type initResp struct{}

// Handle Handle
func (p *Chaos) Handle(c *gin.Context) {
	msgs := make([]define.Msg, 0)
	if err := c.ShouldBind(&msgs); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error()))
		return
	}

	ctx := header.MutateContext(c)
	for _, msg := range msgs {
		if err := p.handler.Put(ctx, msg); err != nil {
			resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error()))
			return
		}
	}
	resp.Format(initResp{}, nil).Context(c)
}
