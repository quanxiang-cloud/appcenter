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

	"github.com/quanxiang-cloud/appcenter/internal/logic"
	"github.com/quanxiang-cloud/appcenter/internal/models"
	"github.com/quanxiang-cloud/appcenter/internal/models/mysql"
	"github.com/quanxiang-cloud/appcenter/internal/req"
	"github.com/quanxiang-cloud/appcenter/internal/resp"
	"github.com/quanxiang-cloud/appcenter/pkg/code"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	page2 "github.com/quanxiang-cloud/appcenter/pkg/page"
	error2 "github.com/quanxiang-cloud/cabin/error"
	id2 "github.com/quanxiang-cloud/cabin/id"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"gorm.io/gorm"
)

type appTemplate struct {
	db           *gorm.DB
	templateRepo models.AppTemplateRepo
	appRepo      models.AppRepo
}

// NewAppTemplate NewAppTemplate
func NewAppTemplate(conf *config.Configs, db *gorm.DB) logic.AppTemplate {
	return &appTemplate{
		db:           db,
		templateRepo: mysql.NewAppTemplateRepo(),
		appRepo:      mysql.NewAppCenterRepo(),
	}
}

func (a *appTemplate) isNameRepeat(ctx context.Context, name string) bool {
	template, err := a.templateRepo.SelectByName(ctx, a.db, name)
	if err != nil || template != nil {
		return true
	}
	return false
}

func (a *appTemplate) preCreate(ctx context.Context, req *req.CreateTemplateReq) error {
	if a.isNameRepeat(ctx, req.Name) {
		return error2.New(code.NameExist)
	}
	return nil
}

func (a *appTemplate) Create(ctx context.Context, req *req.CreateTemplateReq) (*resp.CreateTemplateResp, error) {
	err := a.preCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	appInfo := a.appRepo.SelectByID(req.AppID, a.db)
	if appInfo == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	template := models.AppTemplate{
		ID:          id2.StringUUID(),
		Name:        req.Name,
		AppIcon:     req.AppIcon,
		Path:        req.Path,
		SourceID:    req.AppID,
		SourceName:  appInfo.AppName,
		Version:     req.Version,
		GroupID:     req.GroupID,
		CreatedBy:   req.UserID,
		CreatedName: req.UserName,
		CreatedTime: time2.NowUnix(),
		UpdatedBy:   req.UserID,
		UpdatedName: req.UserName,
		UpdatedTime: time2.NowUnix(),
		Status:      logic.PrivateStatus,
	}
	tx := a.db.Begin()
	err = a.templateRepo.Create(ctx, tx, template)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &resp.CreateTemplateResp{
		ID:      template.ID,
		Name:    template.Name,
		Version: template.Version,
		AppIcon: template.AppIcon,
		AppID:   template.SourceID,
		AppName: template.SourceName,
		GroupID: template.GroupID,
		Status:  template.Status,
	}, nil
}

func (a *appTemplate) Delete(ctx context.Context, req *req.DeleteTemplateReq) (*resp.DeleteTemplateResp, error) {
	template, err := a.templateRepo.SelectByID(ctx, a.db, req.ID)
	if err != nil || template == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	if template.CreatedBy != req.UserID {
		return nil, error2.New(code.ErrNoPermission)
	}

	tx := a.db.Begin()
	err = a.templateRepo.Delete(ctx, tx, req.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &resp.DeleteTemplateResp{}, nil
}

func (a *appTemplate) GetSelfTemplate(ctx context.Context, req *req.GetSelfTemplateReq) (*resp.GetSelfTemplateResp, error) {
	templates, count, err := a.templateRepo.SelectByUser(ctx, a.db, req.Name, req.UserID)
	if err != nil {
		return nil, err
	}
	response := &resp.GetSelfTemplateResp{
		Templates: make([]*resp.TemplateVO, 0, len(templates)),
		Count:     count,
	}
	for _, t := range templates {
		response.Templates = append(response.Templates, templateToVO(t))
	}
	return response, nil
}

func (a *appTemplate) GetTemplatesByPage(ctx context.Context, req *req.GetTemplateByPageReq) (*resp.GetTemplateByPageResp, error) {
	page := page2.NewPage(req.Page, req.PageSize, 0)
	templates, count, err := a.templateRepo.SelectByPage(ctx, a.db, req.Name, models.PublicStatus, page)
	if err != nil {
		return nil, err
	}
	response := &resp.GetTemplateByPageResp{
		Count:     count,
		Page:      page.CurrentPage,
		PageSize:  page.PageSize,
		Templates: make([]*resp.TemplateVO, 0, len(templates)),
	}
	for _, t := range templates {
		response.Templates = append(response.Templates, templateToVO(t))
	}
	return response, nil
}

func (a *appTemplate) GetTemplateByID(ctx context.Context, req *req.GetTemplateByIDReq) (*resp.GetTemplateByIDResp, error) {
	template, err := a.templateRepo.SelectByID(ctx, a.db, req.ID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	return &resp.GetTemplateByIDResp{
		ID:          template.ID,
		Name:        template.Name,
		Version:     template.Version,
		AppIcon:     template.AppIcon,
		Path:        template.Path,
		AppID:       template.SourceID,
		AppName:     template.SourceName,
		GroupID:     template.GroupID,
		CreatedBy:   template.CreatedBy,
		CreatedName: template.CreatedName,
		CreatedTime: template.CreatedTime,
		UpdatedBy:   template.UpdatedBy,
		UpdatedName: template.UpdatedName,
		UpdatedTime: template.UpdatedTime,
		Status:      template.Status,
	}, nil
}

func (a *appTemplate) ModifyStatus(ctx context.Context, req *req.ModifyStatusReq) (*resp.ModifyStatusResp, error) {
	template, err := a.templateRepo.SelectByID(ctx, a.db, req.ID)
	if err != nil || template == nil {
		return nil, error2.New(code.ErrDataNotExist)
	}
	if template.CreatedBy != req.UserID {
		return nil, error2.New(code.ErrNoPermission)
	}
	tx := a.db.Begin()
	err = a.templateRepo.ModifyStatus(ctx, tx, req.ID, req.Status)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &resp.ModifyStatusResp{}, nil
}

func (a *appTemplate) CheckNameRepeat(ctx context.Context, req *req.CheckNameRepeatReq) (*resp.CheckNameRepeatResp, error) {
	return &resp.CheckNameRepeatResp{
		IsRepeat: a.isNameRepeat(ctx, req.Name),
	}, nil
}

func (a *appTemplate) ModifyTemplate(ctx context.Context, req *req.ModifyTemplateReq) (*resp.ModifyTemplateResp, error) {
	template, err := a.templateRepo.SelectByID(ctx, a.db, req.ID)
	if err != nil {
		return nil, err
	}
	if template == nil || template.CreatedBy != req.UserID {
		return nil, error2.New(code.ErrNoPermission)
	}
	if template.Name != req.Name {
		if a.isNameRepeat(ctx, req.Name) {
			return nil, error2.New(code.NameExist)
		}
		template.Name = req.Name
	}
	template.AppIcon = req.AppIcon
	template.UpdatedTime = time2.NowUnix()
	template.UpdatedBy = req.UserID
	template.UpdatedName = req.UserName
	err = a.templateRepo.Update(ctx, a.db, template)
	if err != nil {
		return nil, err
	}
	return &resp.ModifyTemplateResp{
		ID:      template.ID,
		Name:    template.Name,
		AppIcon: template.AppIcon,
	}, nil
}

func (a *appTemplate) FinishCreating(ctx context.Context, req *req.FinishCreatingReq) (*resp.FinishCreatingResp, error) {
	template, err := a.templateRepo.SelectByID(ctx, a.db, req.ID)
	if err != nil {
		return nil, err
	}
	template.Path = req.Path
	template.Status = models.PrivateStatus
	template.UpdatedTime = time2.NowUnix()
	template.UpdatedBy = req.UserID
	template.UpdatedName = req.UserName
	tx := a.db.Begin()
	err = a.templateRepo.Update(ctx, tx, template)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &resp.FinishCreatingResp{}, nil
}

func templateToVO(template models.AppTemplate) *resp.TemplateVO {
	return &resp.TemplateVO{
		ID:          template.ID,
		Name:        template.Name,
		Version:     template.Version,
		AppIcon:     template.AppIcon,
		AppID:       template.SourceID,
		AppName:     template.SourceName,
		GroupID:     template.GroupID,
		CreatedBy:   template.CreatedBy,
		CreatedName: template.CreatedName,
		CreatedTime: template.CreatedTime,
		UpdatedBy:   template.UpdatedBy,
		UpdatedName: template.UpdatedName,
		UpdatedTime: template.UpdatedTime,
		Status:      template.Status,
	}
}
