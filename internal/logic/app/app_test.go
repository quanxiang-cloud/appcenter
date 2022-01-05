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

package app

import (
	"context"

	"github.com/quanxiang-cloud/appcenter/internal/db"
	"github.com/quanxiang-cloud/appcenter/internal/logic"
	"github.com/quanxiang-cloud/appcenter/internal/req"
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/appcenter/pkg/redis"
	"github.com/quanxiang-cloud/cabin/logger"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppCenterSuite struct {
	suite.Suite
	ctx   context.Context
	app   logic.AppCenter
	appID string
}

func _TestAppCenter(t *testing.T) {
	suite.Run(t, new(AppCenterSuite))
}

func (suite *AppCenterSuite) SetupTest() {
	err := config.Init("../../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), config.Config)
	//suite.ctx = logger.ReentryRequsetID(context.Background(), "test-app-center")
	logger.Logger = logger.New(&config.Config.Log)
	assert.Nil(suite.T(), err)
	err = redis.Init()
	assert.Nil(suite.T(), err)
	err = db.InitDB()
	assert.Nil(suite.T(), err)
	suite.app = NewApp(config.Config, db.DB)
}

func (suite *AppCenterSuite) TestSomeAction() {

	rq := req.AddAppCenter{
		AppName:   "test123445",
		AccessURL: "13628005221",
		AppIcon:   "123@test.com",
		CreateBy:  "123",
	}
	res, err := suite.app.Add(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
	suite.appID = res.ID

	rq1 := req.UpdateAppCenter{
		ID:        suite.appID,
		AppName:   "test123445",
		AccessURL: "13628005221",
		AppIcon:   "123@test.com",
		UpdateBy:  "123",
	}
	err1 := suite.app.Update(suite.ctx, &rq1)
	assert.Nil(suite.T(), err1)

	rq2 := req.UpdateAppCenter{
		ID:        suite.appID,
		UseStatus: 1,
		UpdateBy:  "123",
	}
	err2 := suite.app.UpdateStatus(suite.ctx, &rq2)
	assert.Nil(suite.T(), err2)

	suite.UserPageList()

	rq3 := req.SelectOneAppCenter{
		ID: suite.appID,
	}
	res3, err3 := suite.app.AdminSelectByID(suite.ctx, &rq3)
	assert.Nil(suite.T(), err3)
	assert.NotNil(suite.T(), res3)

	rq6 := req.AddAdminUser{
		AppID:   suite.appID,
		UserIDs: []string{"123", "1234"},
	}
	err6 := suite.app.AddAdminUser(suite.ctx, &rq6)
	assert.Nil(suite.T(), err6)

	rq7 := req.SelectAdminUsers{
		ID:    suite.appID,
		Page:  1,
		Limit: 10,
	}
	res7, err7 := suite.app.AdminUsers(suite.ctx, &rq7)
	assert.Nil(suite.T(), res7)
	assert.NotNil(suite.T(), err7)

	rq9 := req.SelectListAppCenter{
		Page:   1,
		Limit:  10,
		UserID: "123",
	}
	res9, err9 := suite.app.AdminPageList(suite.ctx, &rq9)
	assert.Nil(suite.T(), err9)
	assert.NotNil(suite.T(), res9)

	rq10 := req.SelectListAppCenter{
		Page:   1,
		Limit:  10,
		UserID: "123",
	}
	res10, err10 := suite.app.SuperAdminPageList(suite.ctx, &rq10)
	assert.Nil(suite.T(), err10)
	assert.NotNil(suite.T(), res10)

	iDsReq := req.GetAppsByIDsReq{
		IDs: []string{suite.appID},
	}
	res12, err12 := suite.app.GetAppsByIDs(suite.ctx, &iDsReq)
	assert.Nil(suite.T(), err12)
	assert.NotNil(suite.T(), res12)

	isAdminReq := req.CheckIsAdminReq{
		AppID:   suite.appID,
		UserID:  "123",
		IsSuper: false,
	}
	isAdmin := suite.app.CheckIsAdmin(suite.ctx, &isAdminReq)
	assert.NotNil(suite.T(), isAdmin)

	rq8 := req.DelAdminUser{
		AppID:   suite.appID,
		UserIDs: []string{"123", "1234"},
	}
	err8 := suite.app.DelAdminUser(suite.ctx, &rq8)
	assert.Nil(suite.T(), err8)

	rq11 := req.DelAppCenter{
		ID: suite.appID,
	}
	err11 := suite.app.Delete(suite.ctx, &rq11)
	assert.Nil(suite.T(), err11)

}

func (suite *AppCenterSuite) UserPageList() {
	rq := req.SelectListAppCenter{
		Page:  1,
		Limit: 100,
	}
	_, err := suite.app.UserPageList(suite.ctx, &rq)
	assert.Nil(suite.T(), err)
}
