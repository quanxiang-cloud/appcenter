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

package restful

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/appcenter/pkg/broker"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos"
	exec "github.com/quanxiang-cloud/appcenter/pkg/chaos/executor"
	"github.com/quanxiang-cloud/appcenter/pkg/chaos/handle"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/appcenter/pkg/probe"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	ginlog "github.com/quanxiang-cloud/cabin/tailormade/gin"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router router
type Router struct {
	c *config.Configs

	engine *gin.Engine
	Probe  *probe.Probe
}

// NewRouter open the router
func NewRouter(c *config.Configs, log logger.AdaptedLogger) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/app-center")
	db, err := mysql.New(c.Mysql, log)
	if err != nil {
		return nil, err
	}
	app, err := NewAppCenter(c, db)
	if err != nil {
		return nil, err
	}

	k := v1.Group("")
	{
		k.POST("/add", app.Add)
		k.POST("/update", app.Update)
		k.POST("/adminList", checkIsSuperAdmin(app.AdminList, app.SuperAdminList))
		k.POST("/one", app.One)
		k.POST("/addAdmin", app.AddAdmin)
		k.POST("/delAdmin", app.DelAdmin)
		k.POST("/del", app.Del)
		k.POST("/updateStatus", app.UpdateStatus)
		k.POST("/adminUsers", app.AdminUsers)
		k.POST("/checkIsAdmin", app.CheckIsAdmin)
		k.POST("/checkAppAccess", app.CheckAppAccess)

		//----------------------home platform--------------------
		k.POST("/userList", app.UserList)
		k.POST("/apps", app.GetAppsByIDs)

		// -----------------provide services for other services-----------------
		k.POST("/addAppScope", app.AddAppScope)
		k.POST("/getOne", app.GetOne)
		k.POST("/successImport", app.SuccessImport)
		k.POST("/failImport", app.FailImport)
		k.POST("/checkVersion", app.CheckVersion)
		k.POST("/exportApp", app.ExportApp)
		k.POST("/importApp", app.CreateImportApp)
		k.POST("/initCallBack", app.InitCallBack)
	}

	template := NewTemplate(c, db)
	t := v1.Group("/template")
	{
		t.POST("/create", template.Create)
		t.POST("/delete", template.Delete)
		t.POST("/toPublic", template.ToPublic)
		t.POST("/toPrivate", template.ToPrivate)
		t.POST("/publicList", template.GetTemplateByPage)
		t.POST("/selfList", template.GetSelfTemplate)
		t.POST("/getOne", template.GetTemplateByID)
		t.POST("/checkNameRepeat", template.CheckNameRepeat)
		t.POST("/update", template.ModifyTemplate)

		t.POST("/finish", template.FinishCreating)
	}

	r := &Router{
		c:      c,
		engine: engine,
		Probe:  probe.New(),
	}
	r.probe()
	return r, nil
}

// NewInitRouter init router
func NewInitRouter(c *config.Configs, b *broker.Broker, log logger.AdaptedLogger) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	initHandler, err := handle.New(c, b, log)
	if err != nil {
		return nil, err
	}
	// TODO: set executors
	initHandler.SetTaskExecutors(&exec.FormExecutor{
		Client:     client.New(c.InternalNet),
		CreateRole: c.KV[exec.FormCreateRole],
		AssignRole: c.KV[exec.FormAssignRole],
	})
	initHandler.SetSuccessExecutors(&exec.SuccessExecutor{
		BaseExecutor: exec.BaseExecutor{
			Client:       client.New(c.InternalNet),
			AppCenterURL: c.KV[exec.AppCenterURL],
		},
	})
	initHandler.SetFailureExecutors(&exec.FailureExecutor{
		BaseExecutor: exec.BaseExecutor{
			Client:       client.New(c.InternalNet),
			AppCenterURL: c.KV[exec.AppCenterURL],
		},
	})

	p := chaos.New(initHandler, log)
	if err != nil {
		return nil, err
	}
	engine.POST("/init", p.Handle)

	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

func newRouter(c *config.Configs) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()
	engine.Use(ginlog.LoggerFunc(), ginlog.LoggerFunc())
	return engine, nil
}

func (r *Router) probe() {
	r.engine.GET("liveness", func(c *gin.Context) {
		r.Probe.LivenessProbe(c.Writer, c.Request)
	})

	r.engine.Any("readiness", func(c *gin.Context) {
		r.Probe.ReadinessProbe(c.Writer, c.Request)
	})
}

// Run start server
func (r *Router) Run() {
	s := &http.Server{
		Addr:              ":" + r.c.HTTPServer.Port,
		Handler:           r.engine,
		ReadHeaderTimeout: r.c.HTTPServer.ReadHeaderTimeOut * time.Second,
		WriteTimeout:      r.c.HTTPServer.WriteTimeOut * time.Second,
		MaxHeaderBytes:    r.c.HTTPServer.MaxHeaderBytes,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// Close close server
func (r *Router) Close() {
}

func checkIsSuperAdmin(funcAdmin, funcSuperAdmin func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		if isSuperRole(c) {
			funcSuperAdmin(c)
		} else {
			funcAdmin(c)
		}

	}
}
