package chaos

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/handle"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// Chaos Chaos
type Chaos struct {
	log     logger.AdaptedLogger
	handler *handle.InitHandler
}

// New New
func New(handler *handle.InitHandler, log logger.AdaptedLogger, init bool) *Chaos {
	chaos := &Chaos{
		log:     log,
		handler: handler,
	}
	chaos.handler.Run()
	fmt.Println(init)

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
