package exec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
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
	polyReq := &initPolyReq{}
	polyResp := &define.Response{}
	if err := client.POST(ctx, &p.Client, fmt.Sprintf(p.PolyURL, m.AppID), polyReq, polyResp); err != nil {
		return err
	}
	return nil
}

// Bit Bit
func (*PolyExecutor) Bit() int {
	return define.BitPolyAPI
}
