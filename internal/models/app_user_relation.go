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

package models

import "gorm.io/gorm"

// AppUseRelation AppUseRelation
type AppUseRelation struct {
	UserID string `gorm:"column:user_id;type:varchar(64);" json:"userId"`
	AppID  string `gorm:"column:app_id;type:varchar(64);" json:"appId"`
}

//TableName TableName
func (AppUseRelation) TableName() string {
	return "t_app_user_relation"
}

// AppUserRelationRepo AppUserRelationRepo
type AppUserRelationRepo interface {
	Add(rq *AppUseRelation, tx *gorm.DB) (err error)
	DeleteByUserIDAndAppID(appID string, userIDs []string, tx *gorm.DB) (err error)
	DeleteByAppID(appID string, tx *gorm.DB) (err error)
	SelectByAppID(appID string, tx *gorm.DB) (list []AppUseRelation)
	CountByAppIDAndUserID(appID, userID string, db *gorm.DB) int64
	SelectByAppIDBPage(appID string, page, limit int, tx *gorm.DB) (list []AppUseRelation, total int64)
}
