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

package resp

// TemplateVO Template view object
type TemplateVO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppIcon     string `json:"appIcon"`
	AppID       string `json:"appID"`
	AppName     string `json:"appName"`
	GroupID     string `json:"groupID"`
	CreatedBy   string `json:"createdBy"`
	CreatedName string `json:"createdName"`
	CreatedTime int64  `json:"createdTime"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedName string `json:"updatedName"`
	UpdatedTime int64  `json:"updatedTime"`
	Status      int    `json:"status"`
}

// CreateTemplateResp CreateTemplateResp
type CreateTemplateResp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	AppIcon string `json:"appIcon"`
	AppID   string `json:"appID"`
	AppName string `json:"appName"`
	GroupID string `json:"groupID"`
	Status  int    `json:"status"`
}

// DeleteTemplateResp DeleteTemplateResp
type DeleteTemplateResp struct {
}

// GetSelfTemplateResp GetSelfTemplateResp
type GetSelfTemplateResp struct {
	Templates []*TemplateVO `json:"templates"`
	Count     int64         `json:"count"`
}

// GetTemplateByPageResp GetTemplateByPageResp
type GetTemplateByPageResp struct {
	Templates []*TemplateVO `json:"templates"`
	Count     int64         `json:"count"`
	Page      int           `json:"page"`
	PageSize  int           `json:"pageSize"`
}

// GetTemplateByIDResp GetTemplateByIDResp
type GetTemplateByIDResp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppIcon     string `json:"appIcon"`
	Path        string `json:"path"`
	AppID       string `json:"appID"`
	AppName     string `json:"appName"`
	GroupID     string `json:"groupID"`
	CreatedBy   string `json:"createdBy"`
	CreatedName string `json:"createdName"`
	CreatedTime int64  `json:"createdTime"`
	UpdatedBy   string `json:"updatedBy"`
	UpdatedName string `json:"updatedName"`
	UpdatedTime int64  `json:"updatedTime"`
	Status      int    `json:"status"`
}

// ModifyStatusResp ModifyStatusResp
type ModifyStatusResp struct {
}

// CheckNameRepeatResp CheckNameRepeatResp
type CheckNameRepeatResp struct {
	IsRepeat bool `json:"isRepeat"`
}

// ModifyTemplateResp ModifyTemplateResp
type ModifyTemplateResp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	AppIcon string `json:"appIcon"`
}

// FinishCreatingResp FinishCreatingResp
type FinishCreatingResp struct {
}
