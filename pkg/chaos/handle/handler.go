package handle

import (
	"fmt"
	"time"

	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/cabin/logger"
)

type handler func(define.Msg, int) error

func buildExec(executors []Executor) handler {
	return func(msg define.Msg, maximum int) (err error) {
		var count = 0

	RETRY:
		for count < maximum {
			for _, e := range executors {
				if (msg.Content & e.Bit()) == e.Bit() {
					if err = e.Exec(msg); err != nil {
						count++
						continue RETRY
					}
				}
			}
			return
		}
		return
	}
}

type InitHandler struct {
	Stopped bool

	task          chan define.Msg
	taskHandler   handler
	resultHandler handler
	workload      int
	maximumRetry  int

	broker *broker.Broker

	log logger.AdaptedLogger
}

func New(workload, maximumRetry int, broker *broker.Broker, log logger.AdaptedLogger) *InitHandler {
	return &InitHandler{
		task:         make(chan define.Msg, workload*8),
		Stopped:      false,
		broker:       broker,
		log:          log,
		workload:     workload,
		maximumRetry: maximumRetry,
	}
}

func (ih *InitHandler) Put(msg define.Msg) error {
	if !ih.Stopped {
		ih.task <- msg
		return nil
	}
	return fmt.Errorf("handler is stopping")
}

func (ih *InitHandler) Run() {
	if ih.taskHandler == nil {
		ih.log.Warnf("[TaskHandler] taskHandler is a empty func")
		ih.taskHandler = func(define.Msg, int) error {
			ih.log.Warnf("[TaskHandler] empty func is called")
			return nil
		}
	}
	if ih.resultHandler == nil {
		ih.log.Warnf("[ResultHandler] resultHandler is a empty func")
		ih.resultHandler = func(define.Msg, int) error {
			ih.log.Warnf("[ResultHandler] empty func is called")
			return nil
		}
	}
	for i := 0; i < ih.workload; i++ {
		go ih.run()
	}
	ih.withCancel()
}

func (ih *InitHandler) SetTaskExecutors(executors ...Executor) {
	ih.taskHandler = buildExec(executors)
}

func (ih *InitHandler) SetResultExecutors(executors ...Executor) {
	ih.resultHandler = buildExec(executors)
}

func (ih *InitHandler) withCancel() {
	go func() {
		<-ih.broker.C
		ih.Stopped = true

		for len(ih.task) != 0 {
			<-time.After(time.Second)
		}
		ih.broker.Done()
	}()
}

func (ih *InitHandler) run() {
	for {
		msg := <-ih.task
		if err := ih.taskHandler(msg, ih.maximumRetry); err != nil {
			ih.log.Errorf("[TaskHandler] failed to init-server: %s", err.Error())
			continue
		}

		if err := ih.resultHandler(msg, ih.maximumRetry); err != nil {
			ih.log.Errorf("[resultHandler] failed to call-back: %s", err.Error())
		}
	}
}
