package handle

import "github.com/quanxiang-cloud/appcenter/pkg/chaos/define"

type Executor interface {
	Exec(define.Msg) error
	Bit() int
}
