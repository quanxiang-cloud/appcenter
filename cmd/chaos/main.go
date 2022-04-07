package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/appcenter/api/restful"
	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/logger"
)

var (
	configPath = flag.String("config", "../configs/config.yaml", "-config 配置文件地址")
)

func main() {
	flag.Parse()

	config, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
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
