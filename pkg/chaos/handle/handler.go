package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/logger"
)

type data struct {
	Msg   define.Msg      `json:"msg"`
	Ctx   context.Context `json:"ctx"`
	Retry int             `json:"retry"`
	Time  int64           `json:"time"`
}

type handler func(context.Context, define.Msg) (int, error)

func buildExec(executors []Executor) handler {
	return func(ctx context.Context, msg define.Msg) (int, error) {
		for _, e := range executors {
			if (msg.Content&e.Bit()) == e.Bit() && (msg.Ret&e.Bit() == 0) {
				if err := e.Exec(ctx, msg); err != nil {
					return msg.Ret, err
				}
				msg.Ret += e.Bit()
			}
		}
		return msg.Ret, nil
	}
}

// InitHandler InitHandler
type InitHandler struct {
	Stopped bool

	task           chan data
	taskQueue      *taskQueue
	taskHandler    handler
	successHandler handler
	failureHandler handler
	workload       int
	maximumRetry   int
	waitTime       int

	broker *broker.Broker
	log    logger.AdaptedLogger
}

// New New
func New(c *config.Configs, broker *broker.Broker, log logger.AdaptedLogger) (*InitHandler, bool, error) {
	taskQueue, init, err := newTaskQueue(c.CachePath)
	if err != nil {
		return nil, init, err
	}

	handler := &InitHandler{
		task: make(chan data, c.WorkLoad*4), Stopped: false,
		taskQueue:    taskQueue,
		workload:     c.WorkLoad,
		maximumRetry: c.MaximumRetry,
		waitTime:     c.WaitTime,
		broker:       broker,
		log:          log,
	}

	return handler, init, nil
}

// Put Put
func (ih *InitHandler) Put(ctx context.Context, msg define.Msg) error {
	if !ih.Stopped {
		if err := ih.taskQueue.put(data{
			Msg:   msg,
			Ctx:   ctx,
			Retry: 0,
			Time:  time.Now().Unix(),
		}); err != nil {
			return err
		}
	}
	return fmt.Errorf("handler is stopping")
}

// Run Run
func (ih *InitHandler) Run() {
	if ih.taskHandler == nil {
		ih.log.Warnf("[TaskHandler] taskHandler is a empty func")
		ih.taskHandler = func(context.Context, define.Msg) (int, error) {
			ih.log.Warnf("[TaskHandler] empty func is called")
			return 0, nil
		}
	}
	if ih.successHandler == nil {
		ih.log.Warnf("[ResultHandler] resultHandler is a empty func")
		ih.successHandler = func(context.Context, define.Msg) (int, error) {
			ih.log.Warnf("[ResultHandler] empty func is called")
			return 0, nil
		}
	}

	go ih.getTasks()

	for i := 0; i < ih.workload; i++ {
		go ih.run()
	}

	if ih.broker != nil {
		ih.withCancel()
	}
}

func (ih *InitHandler) getTasks() {
	for {
		d, err := ih.taskQueue.pop(ih.workload * 8)
		if err != nil {
			ih.log.Errorf(err.Error())
		}

		for _, one := range d {
			task := &data{}
			if err := json.Unmarshal(one, task); err != nil {
				ih.log.Errorf(err.Error())
				continue
			}
			ih.task <- *task
		}

		time.Sleep(5 * time.Minute)
	}
}

func (ih *InitHandler) run() {
	for {
		data := <-ih.task

		if data.Time <= time.Now().Unix() {
			ret, err := ih.taskHandler(data.Ctx, data.Msg)
			data.Msg.Ret = ret
			if err != nil {
				ih.log.Errorf("[TaskHandler] failed to init-server: %s", err.Error())

				data.Retry++
				if data.Retry < ih.maximumRetry {
					data.Time = time.Now().Add(time.Duration(data.Retry*ih.waitTime) * time.Minute).Unix()
					if err := ih.taskQueue.put(data); err != nil {
						ih.log.Errorf(err.Error())
					}
				} else {
					if _, err := ih.failureHandler(data.Ctx, data.Msg); err != nil {
						ih.log.Errorf("[failureHandler] failed to do: %s", err.Error())
					}
				}
				continue
			}

			if _, err := ih.successHandler(data.Ctx, data.Msg); err != nil {
				ih.log.Errorf("[resultHandler] failed to do: %s", err.Error())
			}
		} else {
			if err := ih.taskQueue.put(data); err != nil {
				ih.log.Errorf(err.Error())
			}
		}

	}
}

// SetTaskExecutors SetTaskExecutors
func (ih *InitHandler) SetTaskExecutors(executors ...Executor) {
	ih.taskHandler = buildExec(executors)
}

// SetSuccessExecutors SetSuccessExecutors
func (ih *InitHandler) SetSuccessExecutors(executors ...Executor) {
	ih.successHandler = buildExec(executors)
}

// SetFailureExecutors SetFailureExecutors
func (ih *InitHandler) SetFailureExecutors(executors ...Executor) {
	ih.failureHandler = buildExec(executors)
}

func (ih *InitHandler) withCancel() {
	go func() {
		<-ih.broker.C
		ih.Stopped = true

		// <-time.After(3 * time.Second)
		for len(ih.task) != 0 {
			<-time.After(time.Second)
		}
		ih.broker.Done()
	}()
}
