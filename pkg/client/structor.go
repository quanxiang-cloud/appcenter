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
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/appcenter/pkg/config"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	structorHost = "/api/v1/structor/%s/base/recycle/app"

	removeTable = structorHost + "/removeTable"

	removePer = structorHost + "/removePer"
)

// NewStructor NewStructor
func NewStructor(conf *config.Configs) Structor {
	return &structor{
		client:     client.New(conf.InternalNet),
		innerHosts: conf.InnerHost,
	}
}

// Structor Structor
type Structor interface {
	RemoveTable(ctx context.Context, appID string) (*DelResp, error)
	RemovePer(ctx context.Context, appID string) (*DelResp, error)
}
type structor struct {
	client     http.Client
	innerHosts config.InnerHostConfig
}

func (s *structor) RemoveTable(ctx context.Context, appID string) (*DelResp, error) {
	params := struct {
		AppID string `json:"appID"`
	}{
		AppID: appID,
	}
	resp := &DelResp{}
	err := client.POST(ctx, &s.client, s.innerHosts.StructorHost+fmt.Sprintf(removeTable, appID), params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

func (s structor) RemovePer(ctx context.Context, appID string) (*DelResp, error) {
	params := struct {
		AppID string `json:"appID"`
	}{
		AppID: appID,
	}
	resp := &DelResp{}
	err := client.POST(ctx, &s.client, s.innerHosts.StructorHost+fmt.Sprintf(removePer, appID), params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
