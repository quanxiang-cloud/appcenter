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

package client

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	appCenterHost = "http://app-center/api/v1/app-center"

	checkIsAdmin = "/checkIsAdmin"
)

// CheckAppAdmin CheckAppAdmin
type CheckAppAdmin struct {
	IsAdmin bool
}

type appCenter struct {
	client http.Client
}

// NewAppCenter NewAppCenter
func NewAppCenter(conf client.Config) AppCenter {
	return &appCenter{
		client: client.New(conf),
	}
}

// AppCenter AppCenter
type AppCenter interface {
	CheckIsAdmin(ctx context.Context, appID, userID string) (CheckAppAdmin, error)
}

func (a *appCenter) CheckIsAdmin(ctx context.Context, appID, userID string) (CheckAppAdmin, error) {
	params := struct {
		AppID  string `json:"appID"`
		UserID string `json:"userID"`
	}{
		AppID:  appID,
		UserID: userID,
	}

	IsAdmin := CheckAppAdmin{}
	err := client.POST(ctx, &a.client, appCenterHost+checkIsAdmin, params, &IsAdmin)
	return IsAdmin, err
}
