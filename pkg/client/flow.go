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

	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	flowHost = "/api/v1/flow/deleteApp"
)

// NewFlow NewFlow
func NewFlow(conf *config.Configs) Flow {
	return &flow{
		client:     client.New(conf.InternalNet),
		innerHosts: conf.InnerHost,
	}
}

// Flow Flow
type Flow interface {
	RemoveApp(ctx context.Context, appID, status string) (*DelResp, error)
}

type flow struct {
	client     http.Client
	innerHosts config.InnerHostConfig
}

// RemoveApp RemoveApp
func (f *flow) RemoveApp(ctx context.Context, appID, status string) (*DelResp, error) {
	params := struct {
		AppID  string `json:"appID"`
		Status string `json:"status"`
	}{
		AppID:  appID,
		Status: status,
	}
	resp := &DelResp{}
	err := client.POST(ctx, &f.client, f.innerHosts.FlowHost+flowHost, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
