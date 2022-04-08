package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/define"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
)

type data struct {
	msg   define.Msg
	ctx   context.Context
	retry int
	time  int64
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
	taskHandler    handler
	successHandler handler
	failureHandler handler
	workload       int
	maximumRetry   int
	waitTime       int
	sync           bool

	redis  *redis.ClusterClient
	broker *broker.Broker
	log    logger.AdaptedLogger
}

// New New
func New(c *config.Configs, broker *broker.Broker, log logger.AdaptedLogger) (*InitHandler, error) {
	handler := &InitHandler{
		task:         make(chan data, c.WorkLoad*8),
		Stopped:      false,
		workload:     c.WorkLoad,
		maximumRetry: c.MaximumRetry,
		waitTime:     c.WaitTime,
		sync:         c.Sync,
		broker:       broker,
		log:          log,
	}

	if handler.sync {
		redis, err := redis2.NewClient(c.Redis)
		if err != nil {
			return nil, err
		}
		handler.redis = redis
	}
	return handler, nil
}

// Put Put
func (ih *InitHandler) Put(ctx context.Context, msg define.Msg) error {
	if !ih.Stopped {
		if ih.sync {
			cache, err := json.Marshal(msg)
			if err != nil {
				return err
			}
			if boolCmd := ih.redis.SetNX(ctx, "chaos:"+msg.AppID, cache, time.Duration(ih.waitTime*ih.maximumRetry)*time.Minute); !boolCmd.Val() {
				ih.log.Warnf("app (%s) is initing.", msg.AppID)
				return nil
			}
		}

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

// Σ(maximumRetry) * waitTime
func cacheDuration(waitTime, maximumRetry int) time.Duration {
	c := ((maximumRetry + 1) * maximumRetry) / 2
	return time.Duration(c*waitTime) * time.Minute
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
	for i := 0; i < ih.workload; i++ {
		go ih.run()
	}

	if ih.broker != nil {
		ih.withCancel()
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
			if ret, err := ih.taskHandler(data.ctx, data.msg); err != nil {
				data.msg.Ret = ret
				ih.log.Errorf("[TaskHandler] failed to init-server: %s", err.Error())

				data.retry++
				if data.retry < ih.maximumRetry {
					data.time = time.Now().Add(time.Duration(data.retry*ih.waitTime) * time.Minute).Unix()
					ih.task <- data
				} else {
					if _, err := ih.failureHandler(data.ctx, data.msg); err != nil {
						ih.log.Errorf("[failureHandler] failed to do: %s", err.Error())
					}
				}
				continue
			}

			if _, err := ih.successHandler(data.ctx, data.msg); err != nil {
				ih.log.Errorf("[resultHandler] failed to do: %s", err.Error())
			}
		} else {
			ih.task <- data
		}
	}
}
