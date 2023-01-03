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

//AppScope AppScope
type AppScope struct {
	AppID   string `gorm:"column:app_id;type:varchar(64)"`
	ScopeID string `gorm:"column:scope_id;type:varchar(64)"`
	Type    string `gorm:"column:type;type:varchar(64)"`
}

// AppUserVO AppUserVO
type AppUserVO struct {
	AppID string `json:"appID"`
	Scope Scope  `json:"scope"`
}

type Scope struct {
	ScopeID string `json:"scopeID"`
	Type    string `json:"type"`
}

// AppScopeRepo AppScopeRepo
type AppScopeRepo interface {
	AppUserDep(db *gorm.DB, appID string, scopes []Scope) error
	DeleteByID(db *gorm.DB, appID string, userID []string) error
	GetByScope(db *gorm.DB, userID, depID string) ([]string, error)
	GetAppByUserID(db *gorm.DB, appID string, userID, depID string) (int64, error)
	GetByAppID(db *gorm.DB, appID string, page, size int) ([]*AppScope, int64, error)
}
