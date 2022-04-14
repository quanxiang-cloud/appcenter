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

type appCenterRepo struct {
}

func (u appCenterRepo) SelectByAppSign(db *gorm.DB, appSign string) *models.AppCenter {
	app := models.AppCenter{}
	affected := db.Where("app_sign=?", appSign).Where("del_flag = 0").Find(&app).RowsAffected
	if affected > 0 {
		return &app
	}
	return nil
}

func (u appCenterRepo) GetDeleteList(db *gorm.DB, deleteTime int64) ([]*models.AppCenter, error) {
	apps := make([]*models.AppCenter, 0)
	err := db.Model(&models.AppCenter{}).
		Where("del_flag = 1 and delete_time <= ? ", deleteTime).
		Find(&apps).
		Error

	return apps, err
}

// UpdateDelFlag mark delete
func (u appCenterRepo) UpdateDelFlag(db *gorm.DB, id string, deleteTime int64) error {
	return db.Model(&models.AppCenter{}).Where("id=?", id).Updates(
		map[string]interface{}{
			"del_flag":    1,
			"delete_time": deleteTime,
		}).Error
}

func (u appCenterRepo) Insert(rq *models.AppCenter, tx *gorm.DB) (err error) {
	err = tx.Create(rq).Error
	if err != nil {
		return err
	}
	return nil
}

func (u appCenterRepo) Update(rq *models.AppCenter, tx *gorm.DB) (err error) {
	err = tx.Model(rq).Updates(rq).Error
	return err
}

func (u appCenterRepo) Delete(id string, tx *gorm.DB) (err error) {
	err = tx.Where("id=?", id).Delete(&models.AppCenter{}).Error
	return err
}

func (u appCenterRepo) SelectByPage(userID, name string, status, page, limit int, isAdmin bool, db *gorm.DB) (list []models.AppCenter, total int64) {
	if name != "" {
		db = db.Where("app_name like ?", "%"+name+"%")
	}
	if isAdmin {
		db = db.Where("id in (select app_id from t_app_user_relation where user_id in (?))", userID)
	}

	if status != 0 {
		db = db.Where("use_status=?", status)
	}
	db = db.Where("del_flag = 0")
	db = db.Order("update_time desc")
	res := make([]models.AppCenter, 0)
	var num int64
	db.Model(&models.AppCenter{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex)

	affected := db.Find(&res).RowsAffected
	if affected > 0 {
		return res, num
	}

	return nil, 0
}

func (u appCenterRepo) SelectByStatus(db *gorm.DB, status int, page, limit int) (list []models.AppCenter, total int64) {
	db = db.Where("use_status=?", status)
	db = db.Where("del_flag = 0")

	res := make([]models.AppCenter, 0)
	var num int64

	db.Model(&models.AppCenter{}).Count(&num)
	newPage := page2.NewPage(page, limit, num)

	db = db.Limit(newPage.PageSize).Offset(newPage.StartIndex)

	affected := db.Find(&res).RowsAffected
	if affected > 0 {
		return res, num
	}
	return nil, 0
}

func (u appCenterRepo) SelectByID(id string, db *gorm.DB) (res *models.AppCenter) {
	app := models.AppCenter{}
	affected := db.Where("id=?", id).Find(&app).RowsAffected
	if affected == 1 {
		return &app
	}
	return nil
}

func (u appCenterRepo) SelectByName(name string, db *gorm.DB) (res *models.AppCenter) {
	app := models.AppCenter{}
	affected := db.Where("app_name=?", name).Where("del_flag = 0").Find(&app).RowsAffected
	if affected > 0 {
		return &app
	}
	return nil
}

func (u appCenterRepo) GetByIDs(tx *gorm.DB, ids ...string) ([]*models.AppCenter, error) {
	apps := make([]*models.AppCenter, 0, len(ids))
	err := tx.Model(&models.AppCenter{}).
		Where("id in ? ", ids).Where("del_flag = 0").Order("update_time desc").
		Find(&apps).
		Error
	return apps, err
}

//NewAppCenterRepo init repo
func NewAppCenterRepo() models.AppRepo {
	return new(appCenterRepo)
}
