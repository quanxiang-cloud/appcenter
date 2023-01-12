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
	"bytes"
	"fmt"

	"github.com/quanxiang-cloud/appcenter/internal/models"
	"gorm.io/gorm"
)

type appScopeRepo struct {
}

func (a *appScopeRepo) GetByScope(db *gorm.DB, userID, depID string) ([]string, error) {
	arr := make([]string, 0)
	if userID != "" {
		arr = append(arr, userID)
	}
	if depID != "" {
		arr = append(arr, depID)
	}
	result := make([]string, 0)
	err := db.Table(a.TableName()).Distinct("app_id").Where("scope_id in  ?", arr).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *appScopeRepo) TableName() string {
	return "t_app_scope"
}
func (a *appScopeRepo) AppUserDep(db *gorm.DB, appID string, scopes []models.Scope) error {
	var buffer bytes.Buffer
	sql := "insert into `t_app_scope` (`app_id`,`scope_id`,`type`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, value := range scopes {
		if i == len(scopes)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s');", appID, value.ScopeID, value.Type))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s'),", appID, value.ScopeID, value.Type))
		}
	}
	return db.Exec(buffer.String()).Error
}

func (a *appScopeRepo) DeleteByID(db *gorm.DB, appID string, scopeIDs []string) error {
	return db.Table(a.TableName()).Where("app_id = ? and scope_id in  ? ", appID, scopeIDs).
		Delete(&models.AppScope{}).
		Error

}

func (a *appScopeRepo) DeleteByAppID(db *gorm.DB, appID string) error {
	return db.Table(a.TableName()).Where("app_id = ?", appID).
		Delete(&models.AppScope{}).
		Error

}

//NewAppScopeRepo init repo
func NewAppScopeRepo() models.AppScopeRepo {
	return &appScopeRepo{}
}

func (a *appScopeRepo) GetAppByUserID(db *gorm.DB, appID string, userID, depID string) (int64, error) {
	arr := make([]string, 0)
	if userID != "" {
		arr = append(arr, userID)
	}
	if depID != "" {
		arr = append(arr, depID)
	}
	ql := db.Table(a.TableName())
	ql = ql.Where("app_id = ? and scope_id in ? ", appID, arr)
	var total int64
	ql.Count(&total)
	return total, nil
}
func (a *appScopeRepo) GetByAppID(db *gorm.DB, appID string, page, size int) ([]*models.AppScope, int64, error) {
	var (
		appScope []*models.AppScope
		count    int64
	)

	ql := db.Table(a.TableName())
	ql = ql.Where("app_id = ?", appID)
	err := ql.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = ql.Offset((page - 1) * size).Limit(size).Find(&appScope).Error
	if err != nil {
		return nil, 0, err
	}

	return appScope, count, nil

}
