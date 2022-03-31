package exec

import (
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
)

type SampleExec struct {
}

func (s *SampleExec) Exec(define.Msg) error {
	// Your logic
	return nil
}

func (*SampleExec) Bit() int {
	// Remake number of bit
	return 1 << 0
}
