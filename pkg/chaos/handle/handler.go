package handle

import (
	"context"
	"fmt"
	"time"

	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/cabin/logger"
)

type data struct {
	msg   define.Msg
	ctx   context.Context
	retry int
	time  int64
}

type handler func(context.Context, define.Msg) error

func buildExec(executors []Executor) handler {
	return func(ctx context.Context, msg define.Msg) error {
		for _, e := range executors {
			if (msg.Content & e.Bit()) == e.Bit() {
				if err := e.Exec(ctx, msg); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

type InitHandler struct {
	Stopped bool

	task           chan data
	taskHandler    handler
	successHandler handler
	failureHandler handler
	workload       int
	maximumRetry   int
	waitTime       int

	broker *broker.Broker

	log logger.AdaptedLogger
}

func New(workload, maximumRetry, waitTime int, broker *broker.Broker, log logger.AdaptedLogger) *InitHandler {
	return &InitHandler{
		task:         make(chan data, workload*8),
		Stopped:      false,
		broker:       broker,
		log:          log,
		workload:     workload,
		maximumRetry: maximumRetry,
		waitTime:     waitTime,
	}
}

func (ih *InitHandler) Put(ctx context.Context, msg define.Msg) error {
	if !ih.Stopped {
		ih.task <- data{
			msg:   msg,
			ctx:   ctx,
			retry: 0,
			time:  time.Now().Unix(),
		}
		return nil
	}
	return fmt.Errorf("handler is stopping")
}

func (ih *InitHandler) Run() {
	if ih.taskHandler == nil {
		ih.log.Warnf("[TaskHandler] taskHandler is a empty func")
		ih.taskHandler = func(context.Context, define.Msg) error {
			ih.log.Warnf("[TaskHandler] empty func is called")
			return nil
		}
	}
	if ih.successHandler == nil {
		ih.log.Warnf("[ResultHandler] resultHandler is a empty func")
		ih.successHandler = func(context.Context, define.Msg) error {
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

func (ih *InitHandler) SetSuccessExecutors(executors ...Executor) {
	ih.successHandler = buildExec(executors)
}

func (ih *InitHandler) SetFailureExecutors(executors ...Executor) {
	ih.failureHandler = buildExec(executors)
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
		data := <-ih.task

		if data.time <= time.Now().Unix() {
			if err := ih.taskHandler(data.ctx, data.msg); err != nil {
				ih.log.Errorf("[TaskHandler] failed to init-server: %s", err.Error())

				data.retry += 1
				if data.retry < ih.maximumRetry {
					data.time = time.Now().Add(time.Duration(data.retry*ih.waitTime) * time.Minute).Unix()
					ih.task <- data
				} else {
					if err := ih.failureHandler(data.ctx, data.msg); err != nil {
						ih.log.Errorf("[failureHandler] failed to do: %s", err.Error())
					}
				}
				continue
			}

			if err := ih.successHandler(data.ctx, data.msg); err != nil {
				ih.log.Errorf("[resultHandler] failed to do: %s", err.Error())
			}
		} else {
			ih.task <- data
		}
	}
}
