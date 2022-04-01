package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quanxiang-cloud/appcenter/api/restful"
	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/cabinet"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

var (
	model             = flag.String("model", "debug", "Server running model. debug | release")
	port              = flag.String("port", "6666", "Port of server.")
	readHeaderTimeOut = flag.Int("readHeadTimeOut", 0, "ReadHeaderTimeout is the amount of time allowed to read request headers.")
	writeTimeout      = flag.Int("writeTimeOut", 0, "WriteTimeout is the maximum duration before timing out writes of the response.")
	maxHeaderBytes    = flag.Int("maxHeaderBytes", 0, "MaxHeaderBytes controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line.")

	logLevel = flag.Int("log-level", -1, "A Level is a logging priority. Higher levels are more important.")

	workLoad     = flag.Int("workLoad", 8, "WorkLoad is the amount of goroutine notifying other server to init.")
	maximumRetry = flag.Int("maximum-retry", 3, "MaximumRetry is the amount of retrying to call init func.")
	waitTime     = flag.Int("waitTime", 2, "WaitTime is the duration of retrying to do task.")

	clientTimeout = flag.Int("clientTimeout", 20, "ClientTimeout is the deadline when dialing other server.")
	maxIdleConns  = flag.Int("maxIdleConns", 10, "MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts.")
)

func main() {
	kv := cabinet.New()
	flag.Parse()

	config := &config.Configs{
		Model: *model,
		HTTPServer: config.HTTPServer{
			Port:              *port,
			ReadHeaderTimeOut: time.Duration(*readHeaderTimeOut),
			WriteTimeOut:      time.Duration(*writeTimeout),
			MaxHeaderBytes:    *maxHeaderBytes,
		},
		Log: logger.Config{
			Level: *logLevel,
		},
		InternalNet: client.Config{
			Timeout:      time.Duration(*clientTimeout),
			MaxIdleConns: *maxIdleConns,
		},
		WorkLoad:     *workLoad,
		MaximumRetry: *maximumRetry,
		WaitTime:     *waitTime,
		KV:           kv,
	}

	log := logger.New(&config.Log)

	broker := broker.New()

	router, err := restful.NewInitRouter(config, broker, log)
	if err != nil {
		panic(err)
	}
	go router.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			router.Close()
			logger.Logger.Sync()
			broker.Cancel()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
