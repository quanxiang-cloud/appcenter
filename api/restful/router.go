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
	"github.com/quanxiang-cloud/appcenter/internal/db"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/appcenter/pkg/probe"
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
func NewRouter(c *config.Configs) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/app-center")
	err = db.InitDB()
	if err != nil {
		return nil, err
	}
	app := NewAppCenter(c, db.DB)

	k := v1.Group("")
	{
		k.POST("/add", app.Add)                                                    //ok
		k.POST("/update", app.Update)                                              //ok
		k.POST("/adminList", checkIsSuperAdmin(app.AdminList, app.SuperAdminList)) //ok
		k.POST("/one", app.One)                                                    //ok
		k.POST("/addAdmin", app.AddAdmin)                                          //ok
		k.POST("/delAdmin", app.DelAdmin)                                          //ok
		k.POST("/del", app.Del)                                                    //ok
		k.POST("/updateStatus", app.UpdateStatus)                                  //ok
		k.POST("/adminUsers", app.AdminUsers)                                      //ok
		k.POST("/checkIsAdmin", app.CheckIsAdmin)                                  //ok
		k.POST("/checkAppAccess", app.CheckAppAccess)
		k.POST("/importApp", app.CreateImportApp)

		//----------------------home platform--------------------
		k.POST("/userList", app.UserList) //ok
		k.POST("/apps", app.GetAppsByIDs)

		// -----------------provide services for other services-----------------
		k.POST("/addAppScope", app.AddAppScope)
		k.POST("/getOne", app.GetOne)
		k.POST("/successImport", app.SuccessImport)
		k.POST("/failImport", app.FailImport)
		k.POST("/checkVersion", app.CheckVersion)
		k.POST("/exportApp", app.ExportApp)
	}

	r := &Router{
		c:      c,
		engine: engine,
		Probe:  probe.New(),
	}
	r.probe()
	return r, nil
}

func newRouter(c *config.Configs) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()
	engine.Use(ginlog.GinLogger(), ginlog.GinRecovery())
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
	s.ListenAndServe()
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
