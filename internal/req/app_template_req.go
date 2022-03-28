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

// CreateTemplateReq CreateTemplateReq
type CreateTemplateReq struct {
	Name     string `json:"name"`
	AppIcon  string `json:"appIcon"`
	AppID    string `json:"appID"`
	Version  string `json:"version"`
	GroupID  string `json:"groupID"`
	Path     string `json:"path"`
	UserID   string `json:"-"`
	UserName string `json:"-"`
}

// DeleteTemplateReq DeleteTemplateReq
type DeleteTemplateReq struct {
	ID     string `json:"id"`
	UserID string `json:"-"`
}

// GetSelfTemplateReq GetSelfTemplateReq
type GetSelfTemplateReq struct {
	Name   string `json:"name"`
	UserID string `json:"-"`
}

// GetTemplateByPageReq GetTemplateByPageReq
type GetTemplateByPageReq struct {
	Name     string `json:"name"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// GetTemplateByIDReq GetTemplateByIDReq
type GetTemplateByIDReq struct {
	ID string `json:"id"`
}

// ModifyStatusReq ModifyStatusReq
type ModifyStatusReq struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	UserID string `json:"-"`
}

// CheckNameRepeatReq CheckNameRepeatReq
type CheckNameRepeatReq struct {
	Name string `json:"name"`
}

// ModifyTemplateReq ModifyTemplateReq
type ModifyTemplateReq struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AppIcon  string `json:"appIcon"`
	UserID   string `json:"-"`
	UserName string `json:"-"`
}

// FinishCreatingReq FinishCreatingReq
type FinishCreatingReq struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	UserID   string `json:"-"`
	UserName string `json:"-"`
}
