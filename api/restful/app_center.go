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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/appcenter/internal/logic"
	"github.com/quanxiang-cloud/appcenter/internal/logic/app"
	"github.com/quanxiang-cloud/appcenter/internal/req"
	resp2 "github.com/quanxiang-cloud/appcenter/internal/resp"
	config2 "github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/logger"
	header2 "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"gorm.io/gorm"
)

const (
	admin  = "X-App-Admin"
	access = "X-App-Access"
)

const (
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
	_roleName     = "Role"
	supperRole    = "super"
)

// AppCenter appCenter
type AppCenter struct {
	appCenter logic.AppCenter
}

// NewAppCenter new appCenter
func NewAppCenter(c *config2.Configs, db *gorm.DB) *AppCenter {
	g := app.NewApp(c, db)
	return &AppCenter{
		appCenter: g,
	}
}

// Add create a app
func (a *AppCenter) Add(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.AddAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	rq.CreateBy = c.GetHeader(_userID)
	rq.CreateByName = c.GetHeader(_userName)
	resp.Format(a.appCenter.Add(ctx, rq)).Context(c)
}

// Update update the app information
func (a *AppCenter) Update(c *gin.Context) {

	ctx := header2.MutateContext(c)
	rq := req.UpdateAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}

	rq.UpdateBy = c.GetHeader(_userID)
	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.ID,
		UserID:  c.GetHeader(_userID),
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	err = a.appCenter.Update(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(nil, nil).Context(c)

}

// AdminList manager get the app list
func (a *AppCenter) AdminList(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.SelectListAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	rq.UserID = c.GetHeader(_userID)

	res, err := a.appCenter.AdminPageList(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)

}

// SuperAdminList super manger get the app list
func (a *AppCenter) SuperAdminList(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.SelectListAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	rq.UserID = c.GetHeader(_userID)

	res, err := a.appCenter.SuperAdminPageList(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)

}

// One find one app
func (a *AppCenter) One(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.SelectOneAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}

	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.ID,
		UserID:  c.GetHeader(_userID),
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	res, err := a.appCenter.AdminSelectByID(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)
}

// AddAdmin add the app manager
func (a *AppCenter) AddAdmin(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.AddAdminUser{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.AppID,
		UserID:  c.GetHeader(_userID),
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	err = a.appCenter.AddAdminUser(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(nil, nil).Context(c)

}

// DelAdmin remove app manager
func (a *AppCenter) DelAdmin(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.DelAdminUser{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.AppID,
		UserID:  c.GetHeader(_userID),
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	err = a.appCenter.DelAdminUser(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(nil, nil).Context(c)
}

// Del delete the application
func (a *AppCenter) Del(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.DelAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}

	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.ID,
		UserID:  c.GetHeader(_userID),
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	err = a.appCenter.Delete(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(nil, nil).Context(c)
}

// UpdateStatus modify the app status
func (a *AppCenter) UpdateStatus(c *gin.Context) {
	ctx := header2.MutateContext(c)

	rq := req.UpdateAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	userID := c.GetHeader(_userID)
	rq.UpdateBy = userID
	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.ID,
		UserID:  userID,
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	err = a.appCenter.UpdateStatus(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(nil, nil).Context(c)

}

// UserList user get the app list on home platform
func (a *AppCenter) UserList(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.SelectListAppCenter{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	rq.UserID = c.GetHeader(_userID)
	departments := strings.Split(c.GetHeader(_departmentID), ",")
	rq.DepID = departments[len(departments)-1]
	res, err := a.appCenter.UserPageList(ctx, &rq)

	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)

}

// AdminUsers get the admin list
func (a *AppCenter) AdminUsers(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.SelectAdminUsers{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}

	isAdminReq := &req.CheckIsAdminReq{
		AppID:   rq.ID,
		UserID:  c.GetHeader(_userID),
		IsSuper: isSuperRole(c),
	}
	isAdmin := a.appCenter.CheckIsAdmin(ctx, isAdminReq)
	if !isAdmin {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	res, err := a.appCenter.AdminUsers(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)

}

// GetAppsByIDs GetAppsByIDs
func (a *AppCenter) GetAppsByIDs(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := req.GetAppsByIDsReq{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	res, err := a.appCenter.GetAppsByIDs(ctx, &rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	resp.Format(res, nil).Context(c)
}

// CheckIsAdmin CheckIsAdmin
func (a *AppCenter) CheckIsAdmin(c *gin.Context) {
	rq := req.CheckIsAdminReq{}
	err := c.ShouldBind(&rq)
	if err != nil {
		logger.Logger.Error(err)
		resp.Format(nil, err).Context(c)
		return
	}
	isAdmin := a.appCenter.CheckIsAdmin(c, &rq)
	res := resp2.CheckIsAdminResp{
		IsAdmin: isAdmin,
	}
	c.Writer.Header().Set(admin, strconv.FormatBool(isAdmin))
	resp.Format(res, nil).Context(c)

}

// AddAppScope AddAppScope
func (a *AppCenter) AddAppScope(c *gin.Context) {
	ctx := header2.MutateContext(c)
	req := &req.AddAppScopeReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.appCenter.AddAppScope(ctx, req)).Context(c)

}

// GetOne GetOne
func (a *AppCenter) GetOne(c *gin.Context) {
	ctx := header2.MutateContext(c)
	req := &req.GetOneReq{}
	if err := c.ShouldBind(req); err != nil {
		logger.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.appCenter.GetOne(ctx, req)).Context(c)

}

// CheckVersion CheckVersion
func (a *AppCenter) CheckVersion(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.CheckImportVersionReq{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.appCenter.CheckImportVersion(ctx, rq)).Context(c)
}

// CreateImportApp CreateImportApp
func (a *AppCenter) CreateImportApp(c *gin.Context) {
	ctx := header2.MutateContext(c)

	rq := &req.AddAppCenter{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.CreateBy = c.GetHeader(_userID)
	rq.CreateByName = c.GetHeader(_userName)
	resp.Format(a.appCenter.CreateImportApp(ctx, rq)).Context(c)
}

// SuccessImport SuccessImport
func (a *AppCenter) SuccessImport(c *gin.Context) {
	ctx := header2.MutateContext(c)

	rq := &req.FinishImportReq{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UpdateBy = c.GetHeader(_userID)
	rq.UpdateByName = c.GetHeader(_userName)
	resp.Format(a.appCenter.FinishImport(ctx, rq)).Context(c)
}

// FailImport FailImport
func (a *AppCenter) FailImport(c *gin.Context) {
	ctx := header2.MutateContext(c)

	rq := &req.ErrorImportReq{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UpdateBy = c.GetHeader(_userID)
	resp.Format(a.appCenter.ErrorImport(ctx, rq)).Context(c)
}

// ExportApp ExportApp
func (a *AppCenter) ExportApp(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.ExportAppReq{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(a.appCenter.ExportApp(ctx, rq)).Context(c)
}

//CheckAppAccess CheckAppAccess
func (a *AppCenter) CheckAppAccess(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.CheckAppAccessReq{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	appAccess, err := a.appCenter.CheckAppAccess(ctx, rq)
	if err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Writer.Header().Set(access, strconv.FormatBool(appAccess.IsAuthority))
}

func isSuperRole(c *gin.Context) bool {
	roles := strings.Split(c.GetHeader(_roleName), ",")
	for _, role := range roles {
		if role == supperRole {
			return true
		}
	}
	return false
}

// InitCallBack call back
func (a *AppCenter) InitCallBack(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.InitCallBackReq{}
	if err := c.ShouldBind(rq); err != nil {
		logger.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UpdateBy = c.GetHeader(_userID)
	rq.UpdateByName = c.GetHeader(_userName)
	resp.Format(a.appCenter.InitCallBack(ctx, rq)).Context(c)
}
