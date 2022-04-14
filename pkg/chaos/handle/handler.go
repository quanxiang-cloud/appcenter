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
	Msg          define.Msg      `json:"msg"`
	SerializeCTX serializeCTX    `json:"ctx"`
	CTX          context.Context `json:"-"`
	Retry        int             `json:"retry"`
	Time         int64           `json:"time"`
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

// TaskHandler TaskHandler
type TaskHandler struct {
	Stopped bool
	Config  *config.Configs

	task         chan data
	taskQueue    *taskQueue
	workload     int
	maximumRetry int
	waitTime     int
	defaultBit   int
	firstInit    bool

	initHandler    InitExecutor
	taskHandler    handler
	successHandler handler
	failureHandler handler

	broker *broker.Broker
	log    logger.AdaptedLogger
}

// New New
func New(c *config.Configs, broker *broker.Broker, log logger.AdaptedLogger) (*TaskHandler, error) {
	taskQueue, init, err := newTaskQueue(c.CachePath)
	if err != nil {
		return nil, err
	}

	handler := &TaskHandler{
		Config: c,
		task:   make(chan data, c.WorkLoad*4), Stopped: false,
		taskQueue:    taskQueue,
		workload:     c.WorkLoad,
		maximumRetry: c.MaximumRetry,
		waitTime:     c.WaitTime,
		defaultBit:   c.InitServerBits,
		firstInit:    init,
		broker:       broker,
		log:          log,
	}

	return handler, nil
}

// Put Put
func (ih *TaskHandler) Put(ctx context.Context, msg define.Msg) error {
	if !ih.Stopped {
		if msg.Content == 0 {
			msg.Content = ih.defaultBit
		}

		d := data{
			Msg:          msg,
			SerializeCTX: marshalCTXHeader(ctx),
			CTX:          ctx,
			Retry:        0,
			Time:         time.Now().Unix(),
		}

		if len(ih.task) < cap(ih.task) {
			ih.task <- d
			return nil
		}

		if err := ih.taskQueue.put(d); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("handler is stopping")
}

// Run Run
func (ih *TaskHandler) Run() error {
	if ih.firstInit {
		if err := ih.initHandler(ih); err != nil {
			return err
		}
	}

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
	return nil
}

func (ih *TaskHandler) getTasks() {
	for {
		d, err := ih.taskQueue.pop(ih.workload * 8)
		if err != nil {
			ih.log.Infof(err.Error())
		}
		if len(d) == 0 {
			time.Sleep(2 * time.Minute)
			continue
		}

		for _, one := range d {
			task := &data{}
			if err := json.Unmarshal(one, task); err != nil {
				ih.log.Errorf(err.Error())
				continue
			}
			task.CTX = unmarshalCTXHeader(task.SerializeCTX)
			ih.task <- *task
		}
	}
}

func (ih *TaskHandler) run() {
	for {
		data := <-ih.task

		if data.Time <= time.Now().Unix() {
			ret, err := ih.taskHandler(data.CTX, data.Msg)
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
					if _, err := ih.failureHandler(data.CTX, data.Msg); err != nil {
						ih.log.Errorf("[failureHandler] failed to do: %s", err.Error())
					}
				}
				continue
			}

			if _, err := ih.successHandler(data.CTX, data.Msg); err != nil {
				ih.log.Errorf("[resultHandler] failed to do: %s", err.Error())
			}
		} else {
			if err := ih.taskQueue.put(data); err != nil {
				ih.log.Errorf(err.Error())
			}
		}

	}
}

// SetInitExecutors SetInitExecutors
func (ih *TaskHandler) SetInitExecutors(executor InitExecutor) {
	ih.initHandler = executor
}

// SetTaskExecutors SetTaskExecutors
func (ih *TaskHandler) SetTaskExecutors(executors ...Executor) {
	ih.taskHandler = buildExec(executors)
}

// SetSuccessExecutors SetSuccessExecutors
func (ih *TaskHandler) SetSuccessExecutors(executors ...Executor) {
	ih.successHandler = buildExec(executors)
}

// SetFailureExecutors SetFailureExecutors
func (ih *TaskHandler) SetFailureExecutors(executors ...Executor) {
	ih.failureHandler = buildExec(executors)
}

func (ih *TaskHandler) withCancel() {
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
