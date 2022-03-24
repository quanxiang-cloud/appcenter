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

package mysql

import (
	"context"

	"github.com/quanxiang-cloud/appcenter/internal/models"
	page2 "github.com/quanxiang-cloud/appcenter/pkg/page"
	"gorm.io/gorm"
)

type appTemplateRepo struct {
}

//NewAppTemplateRepo 初始化
func NewAppTemplateRepo() models.AppTemplateRepo {
	return &appTemplateRepo{}
}

func (a *appTemplateRepo) TableName() string {
	return "t_app_template"
}

func (a *appTemplateRepo) Create(ctx context.Context, tx *gorm.DB, template models.AppTemplate) error {
	err := tx.Table(a.TableName()).Create(template).Error
	return err
}

func (a *appTemplateRepo) Update(ctx context.Context, tx *gorm.DB, template *models.AppTemplate) error {
	err := tx.Table(a.TableName()).Updates(template).Error
	return err
}

func (a *appTemplateRepo) ModifyStatus(ctx context.Context, tx *gorm.DB, id string, status int) error {
	err := tx.Table(a.TableName()).Where("id = ?", id).Update("status", status).Error
	return err
}

func (a *appTemplateRepo) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	err := tx.Table(a.TableName()).Where("id = ?", id).Delete(&models.AppTemplate{}).Error
	return err
}

func (a *appTemplateRepo) SelectByPage(ctx context.Context, db *gorm.DB, name string, status int, page *page2.Page) ([]models.AppTemplate, int64, error) {
	db = db.Table(a.TableName()).Where("status", status)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	db = db.Offset(page.StartIndex).Limit(page.PageSize)
	err = db.Error
	if err != nil {
		return nil, 0, err
	}
	templates := make([]models.AppTemplate, 0, page.PageSize)
	err = db.Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}
	return templates, count, nil
}

func (a *appTemplateRepo) SelectByID(ctx context.Context, db *gorm.DB, id string) (*models.AppTemplate, error) {
	template := &models.AppTemplate{}
	db = db.Table(a.TableName()).Where("id = ?", id).Find(&template)
	err := db.Error
	if err != nil {
		return nil, err
	}
	if db.RowsAffected <= 0 {
		return nil, nil
	}
	return template, nil
}

func (a *appTemplateRepo) SelectByUser(ctx context.Context, db *gorm.DB, name, userID string) ([]models.AppTemplate, int64, error) {
	var count int64
	templates := make([]models.AppTemplate, 0)
	db = db.Table(a.TableName()).Where("created_by = ? and status != ?", userID, models.CreatingStatus)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}
	return templates, count, nil
}

func (a *appTemplateRepo) SelectByName(ctx context.Context, db *gorm.DB, name string) (*models.AppTemplate, error) {
	db = db.Table(a.TableName()).Where("name=? and status!=?", name, models.CreatingStatus)
	template := &models.AppTemplate{}
	err := db.Find(&template).Error
	if err != nil {
		return nil, err
	}
	if db.RowsAffected <= 0 {
		return nil, nil
	}
	return template, nil
}
