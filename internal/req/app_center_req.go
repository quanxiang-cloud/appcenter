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

package req

//AddAppCenter AddAppCenter
type AddAppCenter struct {
	AppName      string                 `json:"appName" binding:"required,max=80,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	AccessURL    string                 `json:"accessURL"`
	AppIcon      string                 `json:"appIcon"`
	CreateBy     string                 `json:"-"`
	CreateByName string                 `json:"_"`
	AppSign      string                 `json:"appSign" binding:"required,alphanum"`
	Extension    map[string]interface{} `json:"extension"`
	Description  string                 `json:"description"`
}

//UpdateAppCenter UpdateAppCenter
type UpdateAppCenter struct {
	ID          string                 `json:"id" binding:"required,max=64"`
	AppName     string                 `json:"appName" binding:"max=80,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	AccessURL   string                 `json:"accessURL"`
	AppIcon     string                 `json:"appIcon"`
	UseStatus   int                    `json:"useStatus"` //published:1，unpublished:-1
	UpdateBy    string                 `json:"-"`
	AppSign     string                 `json:"appSign"`
	Extension   map[string]interface{} `json:"extension"`
	Description string                 `json:"description"`
}

// DelAppCenter DelAppCenter
type DelAppCenter struct {
	ID string `json:"id" binding:"required,max=64"`
}

// SelectListAppCenter SelectListAppCenter
type SelectListAppCenter struct {
	AppName   string `json:"appName" binding:"max=80,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	UseStatus int    `json:"useStatus"` //published:1，unpublished:-1
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	UserID    string `json:"-"`
	DepID     string `json:"depID"`
}

// SelectOneAppCenter SelectOneAppCenter
type SelectOneAppCenter struct {
	ID string `json:"id" binding:"required,max=64"`
}

// AddAdminUser AddAdminUser
type AddAdminUser struct {
	AppID   string   `json:"appID" binding:"required,max=64"`
	UserIDs []string `json:"userIDs" binding:"required,min=1"`
}

// DelAdminUser DelAdminUser
type DelAdminUser struct {
	AppID   string   `json:"appID" binding:"required"`
	UserIDs []string `json:"userIDs" binding:"min=1"`
}

// SelectAdminUsers SelectAdminUsers
type SelectAdminUsers struct {
	ID    string `json:"id" binding:"required,max=64"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// GetAppsByIDsReq GetAppsByIDsReq
type GetAppsByIDsReq struct {
	IDs []string `json:"ids" binding:"required"`
}

// CheckIsAdminReq CheckIsAdminReq
type CheckIsAdminReq struct {
	AppID   string `json:"appID"`
	UserID  string `json:"userID"`
	IsSuper bool   `json:"is_super"`
}

// AddAppScopeReq AddAppScopeReq
type AddAppScopeReq struct {
	AppID  string   `json:"appID"`
	Scopes []string `json:"scopes"`
}

// GetOneReq GetOneReq
type GetOneReq struct {
	AppID string `json:"appID"`
}

// ExportAppReq ExportAppReq
type ExportAppReq struct {
	AppID string `json:"appID" form:"appID"`
}

// ImportAppReq ImportAppReq
type ImportAppReq struct {
	Bytes    []byte
	UserID   string
	UserName string
}

//CheckAppAccessReq CheckAppAccessReq
type CheckAppAccessReq struct {
	AppID  string `json:"appID"`
	UserID string `json:"userID"`
	DepID  string `json:"depID"`
}

// FinishImportReq FinishImportReq
type FinishImportReq struct {
	AppID        string `json:"appID"`
	UpdateBy     string `json:"-"`
	UpdateByName string `json:"-"`
}

// ErrorImportReq ErrorImportReq
type ErrorImportReq struct {
	AppID    string `json:"appID"`
	UpdateBy string `json:"-"`
}

// CheckImportVersionReq CheckImportVersionReq
type CheckImportVersionReq struct {
	Version string `json:"version"`
}

// InitCallBackReq InitCallBackReq
type InitCallBackReq struct {
	ID           string `json:"id"`
	UpdateBy     string `json:"-"`
	UpdateByName string `json:"-"`
	Status       bool   `json:"status"`
	Ret          int    `json:"ret"`
}

// InitServerReq InitServerReq
type InitServerReq struct {
	ID       string `json:"id"`
	CreateBy string `json:"createBy"`
	Server   int    `json:"server"`
}

type ListAppByStatusReq struct {
	Status int `json:"status"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
}
