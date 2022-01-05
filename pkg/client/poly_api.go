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
	polyapiHost = "/api/v1/polyapi/inner"
	requestPoly = polyapiHost + "/request"
	deleteAPP   = polyapiHost + "/delApp/%s"
	parentName  = "/system/poly/permissionInit.p"
)

// NewPolyAPI NewPolyAPI
func NewPolyAPI(conf *config.Configs) PolyAPI {
	return &polyapi{
		client:     client.New(conf.InternalNet),
		innerHosts: conf.InnerHost,
	}
}

type polyapi struct {
	client     http.Client
	innerHosts config.InnerHostConfig
}

// DelResp DelResp
type DelResp struct {
	Errors []*ErrNode `json:"errors"`
}

// ErrNode ErrNode
type ErrNode struct {
	DB    string `json:"db"`
	Table string `json:"table"`
	SQL   string `json:"sql"`
	Err   error  `json:"err"`
}

func (p *polyapi) DeleteAPP(ctx context.Context, appID string) (*DelResp, error) {
	params := struct {
	}{}
	resp := &DelResp{}
	err := client.POST(ctx, &p.client, p.innerHosts.PolyAPI+fmt.Sprintf(deleteAPP, appID), params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RequestPathResp RequestPathResp
type RequestPathResp struct {
}

// RequestPath RequestPath
func (p *polyapi) RequestPath(ctx context.Context, appID, name, description string, types int64, scopes []*ScopesVO) (*RequestPathResp, error) {
	params := struct {
		AppID       string      `json:"appID"`
		Scopes      []*ScopesVO `json:"scopes"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Types       int64       `json:"types"`
	}{
		AppID:       appID,
		Scopes:      scopes,
		Name:        name,
		Description: description,
		Types:       types,
	}
	resp := &RequestPathResp{}
	err := client.POST(ctx, &p.client, p.innerHosts.PolyAPI+requestPoly+parentName, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ScopesVO ScopesVO
type ScopesVO struct {
	Type int16  `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PolyAPI PolyAPI
type PolyAPI interface {
	RequestPath(ctx context.Context, appID, name, description string, types int64, scopes []*ScopesVO) (*RequestPathResp, error)
	DeleteAPP(ctx context.Context, appID string) (*DelResp, error)
}
