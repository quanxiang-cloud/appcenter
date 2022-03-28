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

package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/appcenter/internal/logic"
	"github.com/quanxiang-cloud/appcenter/internal/logic/app"
	"github.com/quanxiang-cloud/appcenter/internal/req"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/logger"
	header2 "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"gorm.io/gorm"
)

// Template Template
type Template struct {
	template logic.AppTemplate
}

// NewTemplate NewTemplate
func NewTemplate(conf *config.Configs, db *gorm.DB) *Template {
	return &Template{
		template: app.NewAppTemplate(conf, db),
	}
}

// Create create template
func (t *Template) Create(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.CreateTemplateReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UserID = c.GetHeader(_userID)
	rq.UserName = c.GetHeader(_userName)
	resp.Format(t.template.Create(ctx, rq)).Context(c)
}

// Delete delete by id
func (t *Template) Delete(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.DeleteTemplateReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UserID = c.GetHeader(_userID)
	resp.Format(t.template.Delete(ctx, rq)).Context(c)
}

// ToPublic make the template to public
func (t *Template) ToPublic(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.ModifyStatusReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UserID = c.GetHeader(_userID)
	rq.Status = logic.PublicStatus
	resp.Format(t.template.ModifyStatus(ctx, rq)).Context(c)
}

// ToPrivate make the template to private
func (t *Template) ToPrivate(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.ModifyStatusReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UserID = c.GetHeader(_userID)
	rq.Status = logic.PrivateStatus
	resp.Format(t.template.ModifyStatus(ctx, rq)).Context(c)
}

// GetSelfTemplate get user's template
func (t *Template) GetSelfTemplate(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.GetSelfTemplateReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	rq.UserID = c.GetHeader(_userID)
	resp.Format(t.template.GetSelfTemplate(ctx, rq)).Context(c)
}

// GetTemplateByID GetTemplateByID
func (t *Template) GetTemplateByID(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.GetTemplateByIDReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.template.GetTemplateByID(ctx, rq)).Context(c)
}

// GetTemplateByPage GetTemplateByPage
func (t *Template) GetTemplateByPage(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.GetTemplateByPageReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.template.GetTemplatesByPage(ctx, rq)).Context(c)
}

// CheckNameRepeat CheckNameRepeat
func (t *Template) CheckNameRepeat(c *gin.Context) {
	ctx := header2.MutateContext(c)
	rq := &req.CheckNameRepeatReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(t.template.CheckNameRepeat(ctx, rq)).Context(c)
}

// ModifyTemplate ModifyTemplate
func (t *Template) ModifyTemplate(c *gin.Context) {
	rq := &req.ModifyTemplateReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx := header2.MutateContext(c)
	rq.UserID = c.GetHeader(_userID)
	rq.UserName = c.GetHeader(_userName)
	resp.Format(t.template.ModifyTemplate(ctx, rq)).Context(c)
}

// FinishCreating FinishCreating
func (t *Template) FinishCreating(c *gin.Context) {
	rq := &req.FinishCreatingReq{}
	if err := c.ShouldBind(&rq); err != nil {
		logger.Logger.Error(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx := header2.MutateContext(c)
	rq.UserID = c.GetHeader(_userID)
	rq.UserName = c.GetHeader(_userName)
	resp.Format(t.template.FinishCreating(ctx, rq)).Context(c)
}
