package chaos

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/handle"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

type Chaos struct {
	log     logger.AdaptedLogger
	handler *handle.InitHandler
}

func New(handler *handle.InitHandler, log logger.AdaptedLogger) *Chaos {
	handler.Run()
	return &Chaos{
		log:     log,
		handler: handler,
	}
}

func (p *Chaos) Handle(c *gin.Context) {
	msg := define.Msg{}
	if err := c.ShouldBind(&msg); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error()))
		return
	}

	if err := p.handler.Put(msg); err != nil {
		resp.Format(nil, error2.NewErrorWithString(error2.ErrParams, err.Error()))
		return
	}
}
