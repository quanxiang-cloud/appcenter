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

import (
	"context"

	page2 "github.com/quanxiang-cloud/appcenter/pkg/page"
	"gorm.io/gorm"
)

const (
	// CreatingStatus CreatingStatus
	CreatingStatus = -1
	// PrivateStatus PrivateStatus
	PrivateStatus = 0
	// PublicStatus PublicStatus
	PublicStatus = 1
)

// AppTemplate Template models
type AppTemplate struct {
	ID          string `gorm:"column:id;type:varchar(64);primary_key;"`
	Name        string `gorm:"column:name;type:varchar(80);"`
	AppIcon     string `gorm:"column:app_icon;type:text;"`
	Path        string `gorm:"column:path;type:varchar(200);"`
	SourceID    string `gorm:"column:source_id;type:varchar(64);"`
	SourceName  string `gorm:"column:source_name;type:varchar(80);"`
	Version     string `gorm:"column:version;type:varchar(64);"`
	GroupID     string `gorm:"column:group_id;type:varchar(64);"`
	CreatedBy   string `gorm:"column:created_by;type:varchar(64);"`
	CreatedName string `gorm:"column:created_name;type:varchar(64);"`
	CreatedTime int64  `gorm:"column:created_time;type:bigint;"`
	UpdatedBy   string `gorm:"column:updated_by;type:varchar(64);"`
	UpdatedName string `gorm:"column:updated_name;type:varchar(64);"`
	UpdatedTime int64  `gorm:"column:updated_time;type:bigint;"`
	Status      int    `gorm:"column:status;type:int;"`
}

// AppTemplateRepo AppTemplateRepo
type AppTemplateRepo interface {
	Create(ctx context.Context, tx *gorm.DB, template AppTemplate) error
	Update(ctx context.Context, tx *gorm.DB, template *AppTemplate) error
	SelectByPage(ctx context.Context, db *gorm.DB, name string, status int, page *page2.Page) ([]AppTemplate, int64, error)
	SelectByID(ctx context.Context, db *gorm.DB, id string) (*AppTemplate, error)
	SelectByUser(ctx context.Context, db *gorm.DB, name, userID string) ([]AppTemplate, int64, error)
	ModifyStatus(ctx context.Context, tx *gorm.DB, id string, status int) error
	Delete(ctx context.Context, tx *gorm.DB, id string) error
	SelectByName(ctx context.Context, db *gorm.DB, name string) (*AppTemplate, error)
}
