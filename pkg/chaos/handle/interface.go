package handle

import (
	"context"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
)

type Executor interface {
	// Exec is the logic func
	Exec(context.Context, define.Msg) error

	// Bit in msg.Content.
	// Exec is always called when bit = 0.
	Bit() int
}
