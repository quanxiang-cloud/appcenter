package exec

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	APP_CENTER_URL = "appcenterurl"
)

type BaseExecutor struct {
	Client       http.Client
	AppCenterURL string

	status bool
}

func (b *BaseExecutor) Exec(ctx context.Context, m define.Msg) error {
	req := &struct {
		ID     string `json:"id"`
		Status bool   `json:"status"`
	}{
		ID:     m.AppID,
		Status: b.status,
	}

	resp := &define.Response{}
	err := client.POST(ctx, &b.Client, b.AppCenterURL, req, resp)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseExecutor) Bit() int {
	return define.BIT_AWAYS
}

type SuccessExecutor struct {
	BaseExecutor
}

func (s *SuccessExecutor) Exec(ctx context.Context, m define.Msg) error {
	s.status = true
	return s.BaseExecutor.Exec(ctx, m)
}

type FailureExecutor struct {
	BaseExecutor
}

func (f *FailureExecutor) Exec(ctx context.Context, m define.Msg) error {
	f.status = false
	return f.BaseExecutor.Exec(ctx, m)
}
