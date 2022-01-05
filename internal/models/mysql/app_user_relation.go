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
	"github.com/quanxiang-cloud/appcenter/internal/models"
	page2 "github.com/quanxiang-cloud/appcenter/pkg/page"
	"gorm.io/gorm"
)

type appUserRelationRepo struct {
}

func (a appUserRelationRepo) Add(rq *models.AppUseRelation, tx *gorm.DB) (err error) {
	err = tx.Create(&rq).Error
	return err
}

func (a appUserRelationRepo) DeleteByUserIDAndAppID(appID string, userIDs []string, tx *gorm.DB) (err error) {
	err = tx.Where("user_id in(?) and app_id=?", userIDs, appID).Delete(models.AppUseRelation{}).Error
	return err
}
func (a appUserRelationRepo) DeleteByAppID(appID string, tx *gorm.DB) (err error) {
	err = tx.Where("app_id=?", appID).Delete(models.AppUseRelation{}).Error
	return err
}
func (a appUserRelationRepo) SelectByAppID(appID string, db *gorm.DB) (list []models.AppUseRelation) {
	relations := make([]models.AppUseRelation, 0)
	affected := db.Where("app_id=?", appID).Find(&relations).RowsAffected
	if affected > 0 {
		return relations
	}
	return nil
}

func (a appUserRelationRepo) SelectByAppIDBPage(appID string, page, limit int, db *gorm.DB) (list []models.AppUseRelation, total int64) {
	relations := make([]models.AppUseRelation, 0)
	var num int64
	db = db.Where("app_id=?", appID)
	db.Model(&models.AppUseRelation{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	affected := db.
		Limit(newPage.PageSize).Offset(newPage.StartIndex).
		Find(&relations).RowsAffected
	if affected > 0 {
		return relations, num
	}
	return nil, 0
}

func (a appUserRelationRepo) CountByAppIDAndUserID(appID, userID string, db *gorm.DB) int64 {
	var num int64
	db = db.Where("app_id=? and user_id=?", appID, userID)
	db.Model(&models.AppUseRelation{}).Count(&num)
	return num
}

//NewAppUserRelationRepo init repo
func NewAppUserRelationRepo() models.AppUserRelationRepo {
	return new(appUserRelationRepo)
}
