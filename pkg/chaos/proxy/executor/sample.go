package exec

import (
	"fmt"

	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
)

type SampleExec struct {
}

func (s *SampleExec) Exec(define.Msg) error {
	fmt.Printf("bit: %d\n", s.Bit())
	return nil
}

func (*SampleExec) Bit() int {
	return 1 << 0
}

type Sample3Exec struct {
}

func (s *Sample3Exec) Exec(define.Msg) error {
	fmt.Printf("bit: %d\n", s.Bit())
	return nil
}

func (*Sample3Exec) Bit() int {
	return 1 << 1
}

type Sample2Exec struct {
}

func (s *Sample2Exec) Exec(define.Msg) error {
	fmt.Printf("bit: %d\n", s.Bit())
	return nil
}

func (*Sample2Exec) Bit() int {
	return 1 << 2
}
