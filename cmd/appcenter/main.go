/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/appcenter/pkg/redis"

	"github.com/quanxiang-cloud/appcenter/api/restful"
	"github.com/quanxiang-cloud/appcenter/pkg/config"

	"github.com/quanxiang-cloud/cabin/logger"
)

var (
	configPath = flag.String("config", "../../configs/config.yml", "-config 配置文件地址")
)

func main() {
	flag.Parse()

	err := config.Init(*configPath)
	if err != nil {
		panic(err)
	}
	config.Config.Model = config.Config.AppCenter.Model
	config.Config.HTTPServer = config.Config.AppCenter.HTTPServer

	logger.Logger = logger.New(&config.Config.Log)

	err = redis.Init()
	if err != nil {
		panic(err)
	}
	// start router
	router, err := restful.NewRouter(config.Config, logger.Logger)
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
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
