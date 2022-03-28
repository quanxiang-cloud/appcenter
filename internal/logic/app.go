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

package logic

import (
	"context"

	"github.com/quanxiang-cloud/appcenter/internal/req"
	"github.com/quanxiang-cloud/appcenter/internal/resp"
	"github.com/quanxiang-cloud/appcenter/pkg/page"
)

// AppCenter app mangement
type AppCenter interface {
	// AdminPageList manager get the app list
	AdminPageList(ctx context.Context, rq *req.SelectListAppCenter) (*page.Page, error)
	// SuperAdminPageList super manager get the app list
	SuperAdminPageList(ctx context.Context, rq *req.SelectListAppCenter) (*page.Page, error)
	// Add create a new app
	Add(ctx context.Context, rq *req.AddAppCenter) (*resp.AdminAppCenter, error)
	// Update modify app information
	Update(ctx context.Context, rq *req.UpdateAppCenter) error
	// UpdateStatus modify app status
	UpdateStatus(ctx context.Context, rq *req.UpdateAppCenter) error
	// Delete delete app
	Delete(ctx context.Context, rq *req.DelAppCenter) error
	// AdminSelectByID get the app details
	AdminSelectByID(ctx context.Context, rq *req.SelectOneAppCenter) (*resp.AdminAppCenter, error)
	// AddAdminUser add the app manager
	AddAdminUser(ctx context.Context, rq *req.AddAdminUser) error
	// DelAdminUser delete the app manager
	DelAdminUser(ctx context.Context, rq *req.DelAdminUser) error
	// AdminUsers AdminUsers
	AdminUsers(ctx context.Context, rq *req.SelectAdminUsers) (*page.Page, error)

	CheckIsAdmin(ctx context.Context, rq *req.CheckIsAdminReq) bool

	// ------Home platform----------

	// UserPageList user get the app list
	UserPageList(ctx context.Context, rq *req.SelectListAppCenter) (*page.Page, error)
	// GetAppsByIDs batch get the app by app ids
	GetAppsByIDs(ctx context.Context, req *req.GetAppsByIDsReq) (*resp.GetAppsByIDsResp, error)
	// AddAppScope AddAppScope
	AddAppScope(ctx context.Context, req *req.AddAppScopeReq) (*resp.AddAppScopeResp, error)
	// GetOne GetOne
	GetOne(ctx context.Context, req *req.GetOneReq) (*resp.GetOneResp, error)
	// CheckAppAccess CheckAppAccess
	CheckAppAccess(ctx context.Context, rq *req.CheckAppAccessReq) (*resp.CheckAppAccessResp, error)

	ExportApp(ctx context.Context, req *req.ExportAppReq) (*resp.ExportAppResp, error)

	FinishImport(ctx context.Context, rq *req.FinishImportReq) (*resp.FinishImportResp, error)

	CreateImportApp(ctx context.Context, rq *req.AddAppCenter) (*resp.AdminAppCenter, error)

	ErrorImport(ctx context.Context, rq *req.ErrorImportReq) (*resp.ErrorImportResp, error)

	CheckImportVersion(ctx context.Context, rq *req.CheckImportVersionReq) (*resp.CheckImportVersionResp, error)
}
