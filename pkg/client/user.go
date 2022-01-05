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
	host = "/api/v1/org"

	usersInfoURI  = "/usersInfo"
	departmentURI = "/depByIDs"
	userDEPIDURI  = "/otherGetUserList"
)

// User User
type User interface {
	// GetInfo GetInfo
	GetInfo(ctx context.Context, userIDs ...string) ([]UserInfo, error)
	// GetDepartment GetDepartment
	GetDepartment(ctx context.Context, ids ...string) ([]Department, error)
	// GetUsersByDEPID GetUsersByDEPID
	GetUsersByDEPID(ctx context.Context, depID string, includeChildDEPChild, page, limit int) ([]UserInfo, error)
}

type user struct {
	client     http.Client
	innerHosts config.InnerHostConfig
}

// NewUser NewUser
func NewUser(conf *config.Configs) User {
	return &user{
		client:     client.New(conf.InternalNet),
		innerHosts: conf.InnerHost,
	}
}

// UserInfo UserInfo
type UserInfo struct {
	ID          string `json:"id"`
	UserName    string `json:"userName"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	LeaderID    string `json:"leaderID"`
	CompanyID   string `json:"companyID"`
	Avatar      string `json:"avatar"`
	IsDEPLeader int    `json:"isDEPLeader,omitempty"`
	DEP         struct {
		ID                 string `json:"id"`
		DepartmentName     string `json:"departmentName"`
		DepartmentLeaderID string `json:"departmentLeaderID"`
		UseStatus          int    `json:"useStatus"`
		PID                string `json:"pid"`
		SuperPID           string `json:"superID"`
		CompanyID          string `json:"companyID"`
		Grade              int    `json:"grade"`
	} `json:"dep"`
}

func (u *user) GetInfo(ctx context.Context, userIDs ...string) ([]UserInfo, error) {
	params := struct {
		IDS []string `json:"ids"`
	}{
		IDS: userIDs,
	}

	userInfo := make([]UserInfo, 0)
	err := client.POST(ctx, &u.client, u.innerHosts.OrgHost+host+usersInfoURI, params, &userInfo)
	return userInfo, err
}

// Department Department
type Department struct {
	ID                 string `json:"id"`
	DepartmentName     string `json:"departmentName"`
	DepartmentLeaderID string `json:"departmentLeaderID"`
	UseStatus          int    `json:"useStatus"`
	PID                string `json:"pid"`
	SuperPID           string `json:"superID"`
	CompanyID          string `json:"companyID"`
	Grade              int    `json:"grade"`
}

func (u *user) GetDepartment(ctx context.Context, ids ...string) ([]Department, error) {
	params := struct {
		IDS []string `json:"ids"`
	}{
		IDS: ids,
	}

	deparment := make([]Department, 0)
	err := client.POST(ctx, &u.client, u.innerHosts.OrgHost+host+departmentURI, params, &deparment)
	return deparment, err
}

func (u *user) GetUsersByDEPID(ctx context.Context, depID string, includeChildDEPChild, page, limit int) ([]UserInfo, error) {
	params := struct {
		DepID                string `json:"depID"`
		IncludeChildDEPChild int    `json:"includeChildDEPChild"`
		Page                 int    `json:"page"`
		Limit                int    `json:"limit"`
	}{
		DepID:                depID,
		IncludeChildDEPChild: includeChildDEPChild,
		Page:                 page,
		Limit:                limit,
	}

	res := make([]UserInfo, 0)
	err := client.POST(ctx, &u.client, u.innerHosts.OrgHost+host+userDEPIDURI, params, &res)
	return res, err
}
