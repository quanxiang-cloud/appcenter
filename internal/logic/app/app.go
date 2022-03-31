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

package app

import (
	"context"
	"time"

	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	time2 "github.com/quanxiang-cloud/cabin/time"

	"github.com/go-redis/redis/v8"

	"github.com/quanxiang-cloud/appcenter/internal/logic"
	"github.com/quanxiang-cloud/appcenter/internal/models"
	"github.com/quanxiang-cloud/appcenter/internal/models/mysql"
	"github.com/quanxiang-cloud/appcenter/internal/req"
	"github.com/quanxiang-cloud/appcenter/internal/resp"
	"github.com/quanxiang-cloud/appcenter/pkg/client"
	"github.com/quanxiang-cloud/appcenter/pkg/code"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/appcenter/pkg/page"
	redis2 "github.com/quanxiang-cloud/appcenter/pkg/redis"

	"gorm.io/gorm"
)

const (
	releaseStatus     = 1
	unReleaseStatus   = -1
	importingStatus   = -2
	errorImportStatus = -3
	appCenterRedis    = "appCenter:admins:"
	randNumber        = 5
	preDelete         = "preDelete"

	changeAdminKey   = "appCenter:admins:change"
	changeAdminValue = "lockValue"
	lockExpTime      = 2
)

// app app
type app struct {
	DB                *gorm.DB
	app               models.AppRepo
	appUser           models.AppUserRelationRepo
	appScope          models.AppScopeRepo
	org               client.User
	redisClient       *redis.ClusterClient
	polyAPI           client.PolyAPI
	flowAPI           client.Flow
	chaosAPI          client.Chaos
	CompatibleVersion string
}

// NewApp return a app instance
func NewApp(c *config.Configs, db *gorm.DB) logic.AppCenter {
	return &app{
		app:               mysql.NewAppCenterRepo(),
		appUser:           mysql.NewAppUserRelationRepo(),
		appScope:          mysql.NewAppScopeRepo(),
		DB:                db,
		org:               client.NewUser(c.InternalNet),
		polyAPI:           client.NewPolyAPI(c),
		redisClient:       redis2.ClusterClient,
		flowAPI:           client.NewFlow(c),
		chaosAPI:          client.NewChaos(c),
		CompatibleVersion: c.CompatibleVersion,
	}
}
func (a *app) AdminPageList(ctx context.Context, rq *req.SelectListAppCenter) (*page.Page, error) {
	list, total := a.app.SelectByPage(rq.UserID, rq.AppName, rq.UseStatus, rq.Page, rq.Limit, true, a.DB)
	if len(list) > 0 {
		res := make([]resp.AdminAppCenter, 0)
		for k := range list {
			appc := resp.AdminAppCenter{}
			appc.ID = list[k].ID
			appc.AppName = list[k].AppName
			appc.AccessURL = list[k].AccessURL
			appc.AppIcon = list[k].AppIcon
			appc.CreateBy = list[k].CreateBy
			appc.UpdateBy = list[k].UpdateBy
			appc.CreateTime = list[k].CreateTime
			appc.UpdateTime = list[k].UpdateTime
			appc.UseStatus = list[k].UseStatus
			res = append(res, appc)
		}
		page := page.Page{}
		page.Data = res
		page.TotalCount = total
		return &page, nil
	}
	return nil, nil
}

func (a *app) SuperAdminPageList(ctx context.Context, rq *req.SelectListAppCenter) (*page.Page, error) {
	list, total := a.app.SelectByPage(rq.UserID, rq.AppName, rq.UseStatus, rq.Page, rq.Limit, false, a.DB)
	if len(list) > 0 {
		res := make([]resp.AdminAppCenter, 0)
		for k := range list {
			appc := resp.AdminAppCenter{}
			appc.ID = list[k].ID
			appc.AppName = list[k].AppName
			appc.AccessURL = list[k].AccessURL
			appc.AppIcon = list[k].AppIcon
			appc.CreateBy = list[k].CreateBy
			appc.UpdateBy = list[k].UpdateBy
			appc.CreateTime = list[k].CreateTime
			appc.UpdateTime = list[k].UpdateTime
			appc.UseStatus = list[k].UseStatus
			res = append(res, appc)
		}
		page := page.Page{}
		page.Data = res
		page.TotalCount = total
		return &page, nil
	}
	return nil, nil
}

func (a *app) Add(ctx context.Context, rq *req.AddAppCenter) (*resp.AdminAppCenter, error) {
	appCenter := a.app.SelectByName(rq.AppName, a.DB)
	if appCenter != nil {
		return nil, error2.New(code.NameExist)
	}
	appCenter = a.app.SelectByAppSign(a.DB, rq.AppSign)
	if appCenter != nil {
		return nil, error2.New(code.ErrIdentifiesExist)
	}

	app := models.AppCenter{}
	nowUnix := time2.NowUnix()
	id := id2.String(randNumber)
	app.ID = id
	app.AppName = rq.AppName
	app.AccessURL = rq.AccessURL
	app.AppIcon = rq.AppIcon
	app.CreateBy = rq.CreateBy
	app.UpdateBy = rq.CreateBy
	app.CreateTime = nowUnix
	app.UpdateTime = nowUnix
	app.UseStatus = unReleaseStatus
	app.AppSign = rq.AppSign
	tx := a.DB.Begin()
	err := a.app.Insert(&app, tx)
	if err != nil {
		return nil, err
	}
	relation := models.AppUseRelation{}
	relation.UserID = rq.CreateBy
	relation.AppID = id
	err = a.appUser.Add(&relation, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	center := resp.AdminAppCenter{
		ID:       id,
		CreateBy: rq.CreateBy,
	}
	tx.Commit()

	// // init server
	// scopes := make([]*client.ScopesVO, 0)
	// scope := &client.ScopesVO{
	// 	ID:   rq.CreateBy,
	// 	Type: 1,
	// 	Name: rq.CreateByName,
	// }
	// scopes = append(scopes, scope)
	// _, err = a.polyAPI.RequestPath(ctx, id, name, description, perInitTypes, scopes)
	// if err != nil {
	// 	return nil, err
	// }

	// TODO: init content
	a.chaosAPI.Init(ctx, id, rq.CreateBy, rq.CreateByName, 0)

	return &center, nil
}

func (a *app) Update(ctx context.Context, rq *req.UpdateAppCenter) error {
	center := a.app.SelectByID(rq.ID, a.DB)
	if center == nil {
		return error2.New(code.InvalidParams)
	}
	appc := models.AppCenter{}
	if center.AppName != rq.AppName {
		if rq.AppName != "" {
			appCenter := a.app.SelectByName(rq.AppName, a.DB)
			if appCenter != nil {
				return error2.New(code.NameExist)
			}
			appc.AppName = rq.AppName
		}
	}
	if center.AppSign == "" && rq.AppSign != "" {
		ac := a.app.SelectByAppSign(a.DB, rq.AppSign)
		if ac != nil {
			return error2.New(code.ErrIdentifiesExist)
		}
		appc.AppSign = rq.AppSign
	}
	nowUnix := time2.NowUnix()
	appc.ID = rq.ID
	appc.AccessURL = rq.AccessURL
	appc.AppIcon = rq.AppIcon
	appc.UpdateBy = rq.UpdateBy
	appc.UpdateTime = nowUnix
	tx := a.DB.Begin()
	err := a.app.Update(&appc, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (a *app) UpdateStatus(ctx context.Context, rq *req.UpdateAppCenter) error {
	appc := models.AppCenter{}
	nowUnix := time2.NowUnix()
	appc.ID = rq.ID
	appc.UseStatus = rq.UseStatus
	appc.UpdateBy = rq.UpdateBy
	appc.UpdateTime = nowUnix
	tx := a.DB.Begin()
	err := a.app.Update(&appc, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (a *app) Delete(ctx context.Context, rq *req.DelAppCenter) error {
	FiveDayTime := time.Now().AddDate(0, 0, 5)                        // Get the time five days later
	err := a.app.UpdateDelFlag(a.DB, rq.ID, FiveDayTime.UTC().Unix()) // Mark deletion
	if err != nil {
		return err
	}
	_, err = a.flowAPI.RemoveApp(ctx, rq.ID, preDelete)
	if err != nil {
		logger.Logger.Error("delete flow is error ", err.Error())
	}
	return nil
}

func (a *app) AdminSelectByID(ctx context.Context, rq *req.SelectOneAppCenter) (*resp.AdminAppCenter, error) {
	appc := a.app.SelectByID(rq.ID, a.DB)
	if appc != nil {
		res := resp.AdminAppCenter{}
		res.ID = appc.ID
		res.AppName = appc.AppName
		res.AccessURL = appc.AccessURL
		res.AppIcon = appc.AppIcon
		res.CreateBy = appc.CreateBy
		res.UpdateBy = appc.UpdateBy
		res.CreateTime = appc.CreateTime
		res.UpdateTime = appc.UpdateTime
		res.UseStatus = appc.UseStatus
		res.DelFlag = appc.DelFlag
		res.AppSign = appc.AppSign
		return &res, nil
	}
	return nil, nil
}

func (a *app) AddAdminUser(ctx context.Context, rq *req.AddAdminUser) error {
	tx := a.DB.Begin()
	err := a.appUser.DeleteByAppID(rq.AppID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	for k := range rq.UserIDs {
		relation := models.AppUseRelation{}
		relation.AppID = rq.AppID
		relation.UserID = rq.UserIDs[k]
		err := a.appUser.Add(&relation, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	locker := redis2.NewLocker(changeAdminKey, changeAdminValue, lockExpTime, a.redisClient)
	for {
		out := time.After(lockExpTime * 100 * time.Millisecond)

		<-out
		lock, err := locker.Lock()
		if err != nil {
			return err
		}
		if lock {
			err = a.redisAdminUserCacheUpdate(ctx, rq.AppID, rq.UserIDs)
			if err != nil {
				logger.Logger.Error("delete is error", err.Error())
			}
			locker.UnLock()
			return nil
		}

	}
}

func (a *app) DelAdminUser(ctx context.Context, rq *req.DelAdminUser) error {
	if len(rq.UserIDs) > 0 {
		tx := a.DB.Begin()
		err := a.appUser.DeleteByUserIDAndAppID(rq.AppID, rq.UserIDs, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		relations := a.appUser.SelectByAppID(rq.AppID, a.DB)
		userIDs := make([]string, 0)
		for k := range relations {
			userIDs = append(userIDs, relations[k].UserID)
		}
		locker := redis2.NewLocker(changeAdminKey, changeAdminValue, lockExpTime, a.redisClient)
		for {
			out := time.After(lockExpTime * 100 * time.Millisecond)

			<-out
			lock, err := locker.Lock()
			if err != nil {
				return err
			}
			if lock {
				err = a.redisAdminUserCacheUpdate(ctx, rq.AppID, userIDs)
				if err != nil {
					logger.Logger.Error("delete redis is error ", err.Error())
				}
				locker.UnLock()
				return nil
			}

		}
	}
	return error2.New(code.InvalidDel)
}

func (a *app) AdminUsers(ctx context.Context, rq *req.SelectAdminUsers) (*page.Page, error) {
	relations, total := a.appUser.SelectByAppIDBPage(rq.ID, rq.Page, rq.Limit, a.DB)
	p := page.Page{}
	if len(relations) > 0 {
		ids := make([]string, 0)
		for k := range relations {
			ids = append(ids, relations[k].UserID)
		}

		userInfos, err := a.org.GetUserByIDs(ctx, &client.GetUserByIDsRequest{
			IDs: ids,
		})
		if err != nil {
			return nil, err
		}
		p.Data = userInfos
		p.TotalCount = total
		return &p, nil
	}
	return &p, nil
}

//---------------------Home platform------------------------

// UserPageList UserPageList
func (a *app) UserPageList(ctx context.Context, rq *req.SelectListAppCenter) (*page.Page, error) {
	// find appID
	appIDs, err := a.appScope.GetByScope(a.DB, rq.UserID, rq.DepID)
	if err != nil {
		return nil, err
	}
	list, err := a.app.GetByIDs(a.DB, appIDs...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		res := make([]resp.UserAppCenter, 0)
		for k := range list {
			if list[k].UseStatus == releaseStatus {
				appc := resp.UserAppCenter{}
				appc.ID = list[k].ID
				appc.AppName = list[k].AppName
				appc.AccessURL = list[k].AccessURL
				appc.AppIcon = list[k].AppIcon
				res = append(res, appc)
			}
		}
		page := page.Page{}
		page.Data = res
		page.TotalCount = int64(len(res))
		return &page, nil
	}
	return nil, nil
}

// GetAppsByIDs GetAppsByIDs
func (a *app) GetAppsByIDs(ctx context.Context, req *req.GetAppsByIDsReq) (*resp.GetAppsByIDsResp, error) {
	apps, err := a.app.GetByIDs(a.DB, req.IDs...)
	if err != nil {
		return nil, err
	}
	result := &resp.GetAppsByIDsResp{
		Apps: make([]*resp.UserAppCenter, 0, len(apps)),
	}

	for _, appc := range apps {
		result.Apps = append(result.Apps, &resp.UserAppCenter{
			ID:        appc.ID,
			AppName:   appc.AppName,
			AppIcon:   appc.AppIcon,
			AccessURL: appc.AccessURL,
		})
	}

	return result, nil
}

// CheckIsAdmin CheckIsAdmin
func (a *app) CheckIsAdmin(ctx context.Context, rq *req.CheckIsAdminReq) bool {
	app := a.app.SelectByID(rq.AppID, a.DB)
	if app == nil || app.DelFlag == models.Deleted {
		return false
	}
	if !rq.IsSuper {
		val := a.redisClient.HExists(ctx, appCenterRedis+rq.AppID, rq.UserID).Val()
		if val {
			return true
		}
		num := a.appUser.CountByAppIDAndUserID(rq.AppID, rq.UserID, a.DB)
		return num > 0
	}
	return true
}

//  redisAdminUserCacheUpdate  redisAdminUserCacheUpdate
func (a *app) redisAdminUserCacheUpdate(ctx context.Context, appID string, userIDs []string) error {
	usersID := a.redisClient.HKeys(ctx, appCenterRedis+appID).Val()
	if len(usersID) > 0 {
		err := a.redisClient.Del(ctx, appCenterRedis+appID).Err()
		if err != nil {
			return err
		}
	}
	for k := range userIDs {
		err := a.redisClient.HSet(ctx, appCenterRedis+appID, userIDs[k], userIDs[k]).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// AddAppScope AddAppScope
func (a *app) AddAppScope(ctx context.Context, req *req.AddAppScopeReq) (*resp.AddAppScopeResp, error) {
	tx := a.DB.Begin()
	// 1. delete
	err := a.appScope.DeleteByID(tx, req.AppID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = a.appScope.AppUserDep(tx, req.AppID, req.Scopes)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &resp.AddAppScopeResp{}, nil
}

// GetOne GetOne
func (a *app) GetOne(ctx context.Context, req *req.GetOneReq) (*resp.GetOneResp, error) {
	appc := a.app.SelectByID(req.AppID, a.DB)
	if appc != nil {
		res := resp.GetOneResp{}
		res.ID = appc.ID
		res.Name = appc.AppName
		res.DelFlag = appc.DelFlag
		return &res, nil
	}
	return nil, nil
}

func (a *app) ExportApp(ctx context.Context, req *req.ExportAppReq) (*resp.ExportAppResp, error) {
	app := a.app.SelectByID(req.AppID, a.DB)
	if app == nil {
		return nil, error2.New(code.InvalidURI)
	}
	return &resp.ExportAppResp{
		AppID:   app.ID,
		AppName: app.AppName,
		Version: a.CompatibleVersion,
	}, nil
}

func (a *app) CreateImportApp(ctx context.Context, rq *req.AddAppCenter) (*resp.AdminAppCenter, error) {
	appCenter := a.app.SelectByName(rq.AppName, a.DB)
	if appCenter != nil {
		return nil, error2.New(code.NameExist)
	}
	appCenter = a.app.SelectByAppSign(a.DB, rq.AppSign)
	if appCenter != nil {
		return nil, error2.New(code.ErrIdentifiesExist)
	}
	app := models.AppCenter{}
	nowUnix := time2.NowUnix()
	id := id2.String(randNumber)
	app.ID = id
	app.AppName = rq.AppName
	app.AccessURL = rq.AccessURL
	app.AppIcon = rq.AppIcon
	app.CreateBy = rq.CreateBy
	app.UpdateBy = rq.CreateBy
	app.CreateTime = nowUnix
	app.UpdateTime = nowUnix
	app.UseStatus = importingStatus
	app.AppSign = rq.AppSign
	tx := a.DB.Begin()
	err := a.app.Insert(&app, tx)
	if err != nil {
		return nil, err
	}
	relation := models.AppUseRelation{}
	relation.UserID = rq.CreateBy
	relation.AppID = id
	err = a.appUser.Add(&relation, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	center := resp.AdminAppCenter{
		ID:       id,
		CreateBy: rq.CreateBy,
	}
	tx.Commit()
	return &center, nil
}

func (a *app) FinishImport(ctx context.Context, rq *req.FinishImportReq) (*resp.FinishImportResp, error) {
	err := a.UpdateStatus(ctx, &req.UpdateAppCenter{
		ID:        rq.AppID,
		UpdateBy:  rq.UpdateBy,
		UseStatus: unReleaseStatus,
	})
	if err != nil {
		return nil, err
	}
	return &resp.FinishImportResp{}, nil
}

func (a *app) ErrorImport(ctx context.Context, rq *req.ErrorImportReq) (*resp.ErrorImportResp, error) {
	err := a.UpdateStatus(ctx, &req.UpdateAppCenter{
		ID:        rq.AppID,
		UpdateBy:  rq.UpdateBy,
		UseStatus: errorImportStatus,
	})
	if err != nil {
		return nil, err
	}
	return &resp.ErrorImportResp{}, nil
}

func (a *app) CheckImportVersion(ctx context.Context, rq *req.CheckImportVersionReq) (*resp.CheckImportVersionResp, error) {
	if rq.Version != a.CompatibleVersion {
		return nil, error2.New(code.ErrVersion)
	}
	return &resp.CheckImportVersionResp{}, nil
}

func (a *app) CheckAppAccess(ctx context.Context, rq *req.CheckAppAccessReq) (*resp.CheckAppAccessResp, error) {
	app := a.app.SelectByID(rq.AppID, a.DB)
	if app.DelFlag == models.Deleted {
		return &resp.CheckAppAccessResp{
			IsAuthority: false,
		}, nil
	}
	appIDCount, err := a.appScope.GetAppByUserID(a.DB, rq.AppID, rq.UserID, rq.DepID)
	if err != nil {
		return nil, err
	}
	return &resp.CheckAppAccessResp{
		IsAuthority: appIDCount > 0,
	}, nil
}
