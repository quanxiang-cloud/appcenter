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

package logic

import (
	"context"

	"github.com/quanxiang-cloud/appcenter/internal/req"
	"github.com/quanxiang-cloud/appcenter/internal/resp"
)

const (
	// PrivateStatus PrivateStatus
	PrivateStatus = 0
	// PublicStatus PublicStatus
	PublicStatus = 1
)

// AppTemplate AppTemplate
type AppTemplate interface {
	Create(ctx context.Context, req *req.CreateTemplateReq) (*resp.CreateTemplateResp, error)
	Delete(ctx context.Context, req *req.DeleteTemplateReq) (*resp.DeleteTemplateResp, error)
	GetSelfTemplate(ctx context.Context, req *req.GetSelfTemplateReq) (*resp.GetSelfTemplateResp, error)
	GetTemplatesByPage(ctx context.Context, req *req.GetTemplateByPageReq) (*resp.GetTemplateByPageResp, error)
	GetTemplateByID(ctx context.Context, req *req.GetTemplateByIDReq) (*resp.GetTemplateByIDResp, error)
	ModifyStatus(ctx context.Context, req *req.ModifyStatusReq) (*resp.ModifyStatusResp, error)
	CheckNameRepeat(ctx context.Context, req *req.CheckNameRepeatReq) (*resp.CheckNameRepeatResp, error)
	ModifyTemplate(ctx context.Context, req *req.ModifyTemplateReq) (*resp.ModifyTemplateResp, error)
	FinishCreating(ctx context.Context, req *req.FinishCreatingReq) (*resp.FinishCreatingResp, error)
}
