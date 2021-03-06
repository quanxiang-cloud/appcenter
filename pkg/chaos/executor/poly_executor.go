package exec

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/cabin/logger"
)

// Key
const (
	PolyInit = "poly-init"
)

// PolyExecutor PolyExecutor
type PolyExecutor struct {
	Client  http.Client
	PolyURL string
}

type initPolyReq struct {
	Data initAppPath `json:"data"`
}

type initAppPath struct {
	AppID string `json:"appID"`
}

// Exec Exec
func (p *PolyExecutor) Exec(ctx context.Context, m define.Msg) error {
	polyReq := &initPolyReq{
		Data: initAppPath{
			AppID: m.AppID,
		},
	}
	polyResp := &define.Response{}
	if err := post(ctx, &p.Client, p.PolyURL, polyReq, polyResp); err != nil {
		logger.Logger.Errorf("init polyapi url: %s", p.PolyURL)
		logger.Logger.Errorf("%s", err)
		return err
	}
	return nil
}

// Bit Bit
func (*PolyExecutor) Bit() int {
	return define.BitPolyAPI
}
