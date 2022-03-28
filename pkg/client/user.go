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
	host = "http://org/api/v1/org"

	othAddUsersURI  = "/o/user/add"
	othAddDepsURI   = "/o/department/add"
	oneUserURI      = "/o/user/info"
	usersByIDsURI   = "/o/user/ids"
	depByIDsURI     = "/o/dep/ids"
	usersByDepIDURI = "/o/user/dep/id"
	depMaxGradeURI  = "/o/dep/max/grade"
)

// User organization service
type User interface {
	OthAddUsers(ctx context.Context, r *AddUsersRequest) (*AddListResponse, error)
	OthAddDeps(ctx context.Context, r *AddDepartmentRequest) (*AddListResponse, error)
	GetUserInfo(ctx context.Context, r *OneUserRequest) (*OneUserResponse, error)
	GetUserByIDs(ctx context.Context, r *GetUserByIDsRequest) (*GetUserByIDsResponse, error)
	GetDepByIDs(ctx context.Context, r *GetDepByIDsRequest) (*GetDepByIDsResponse, error)
	GetUsersByDepID(ctx context.Context, r *GetUsersByDepIDRequest) (*GetUsersByDepIDResponse, error)
	GetDepMaxGrade(ctx context.Context, r *GetDepMaxGradeRequest) (*GetDepMaxGradeResponse, error)
}
type user struct {
	client http.Client
}

// NewUser init instance
func NewUser(conf client.Config) User {
	return &user{
		client: client.New(conf),
	}
}

//AddUsersRequest other server add user request
type AddUsersRequest struct {
	Users      []AddUser `json:"users"`
	IsUpdate   int       `json:"isUpdate"`   //Whether to update existing data，1:true，-1:not
	SyncID     string    `json:"syncID"`     //Synchronization center ID
	SyncSource string    `json:"syncSource"` //Synchronization source
}

//AddUser other server add user to org
type AddUser struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Email     string   `json:"email,omitempty"`
	AccountID string   `json:"-"`
	SelfEmail string   `json:"selfEmail,omitempty"`
	IDCard    string   `json:"idCard,omitempty"`
	Address   string   `json:"address,omitempty"`
	UseStatus int      `json:"useStatus,omitempty"` //Status: 1 Normal, -2 Disabled, -1 Deleted, 2 Activated ==1, -3 deactivated (same as account library)
	Gender    int      `json:"gender,omitempty"`    //Gender: 0 None, 1 male, 2 female
	CompanyID string   `json:"companyID,omitempty"` //Id of Company
	Position  string   `json:"position,omitempty"`  //Position
	Avatar    string   `json:"avatar,omitempty"`    //Head
	Remark    string   `json:"remark,omitempty"`    //Note
	JobNumber string   `json:"jobNumber,omitempty"` //Work number
	DepIDs    []string `json:"depIDs,omitempty"`
	EntryTime int64    `json:"entryTime,omitempty" ` //Hiredate
	Source    string   `json:"source,omitempty" `    //Source of information
	SourceID  string   `json:"sourceID,omitempty" `  //The ID of the source of the information to be returned to the service
}

//AddListResponse other server add user or dep to org response
type AddListResponse struct {
	Result map[int]*Result `json:"result"`
}

//Result list add response
type Result struct {
	ID     string `json:"id"`
	Remark string `json:"remark"`
	Attr   int    `json:"attr"` //11 add ok,0fail,12, update ok
}

//OthAddUsers OthAddUsers
func (u *user) OthAddUsers(ctx context.Context, r *AddUsersRequest) (*AddListResponse, error) {
	response := &AddListResponse{}
	err := client.POST(ctx, &u.client, host+othAddUsersURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

//AddDepartmentRequest other server add  department to org request
type AddDepartmentRequest struct {
	Deps       []AddDep `json:"deps"`
	SyncDep    int      `json:"syncDep"`    //Department 1 is synchronized. -1 is not synchronized
	IsUpdate   int      `json:"isUpdate"`   //Whether to update existing data. 1 is updated and -1 is not updated
	SyncID     string   `json:"syncID"`     //Synchronization center ID
	SyncSource string   `json:"syncSource"` //Synchronization source
}

//AddDep other server add department to org
type AddDep struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UseStatus int    `json:"useStatus"` //1 Normal, -1 true deletion, -2 disabled
	Attr      int    `json:"attr"`      //1 company, 2 department
	PID       string `json:"pid"`       //The upper ID
	SuperPID  string `json:"superID"`   //ID of the top-level parent
	CompanyID string `json:"companyID"` //Id of Company
	Grade     int    `json:"grade"`     //Department level
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	CreatedBy string `json:"createdBy"`        //The creator
	UpdatedBy string `json:"updatedBy"`        //The modifier
	Remark    string `json:"remark,omitempty"` //Note
}

//OthAddDeps The actual request
func (u *user) OthAddDeps(ctx context.Context, r *AddDepartmentRequest) (*AddListResponse, error) {
	response := &AddListResponse{}
	err := client.POST(ctx, &u.client, host+othAddDepsURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// OneUserRequest Query one
type OneUserRequest struct {
	ID string `json:"id" form:"id"  binding:"required,max=64"`
}

// OneUserResponse OneUserResponse
type OneUserResponse struct {
	ID        string              `json:"id,omitempty" `
	Name      string              `json:"name,omitempty" `
	Phone     string              `json:"phone,omitempty" `
	Email     string              `json:"email,omitempty" `
	SelfEmail string              `json:"selfEmail,omitempty" `
	UseStatus int                 `json:"useStatus,omitempty" ` // Status: 1 Normal, -2 disabled, -3 demission, -1 Deleted, 2 Active ==1 (same as account library)
	TenantID  string              `json:"tenantID,omitempty" `  // the tenant id
	Position  string              `json:"position,omitempty" `  // position
	Avatar    string              `json:"avatar,omitempty" `    //
	JobNumber string              `json:"jobNumber,omitempty" `
	Status    int                 `json:"status"`
	Dep       [][]DepOneResponse  `json:"deps,omitempty"`
	Leader    [][]OneUserResponse `json:"leaders,omitempty"`
}

// DepOneResponse DepOneResponse
type DepOneResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	LeaderID  string `json:"leaderID"`
	UseStatus int    `json:"useStatus,omitempty"`
	PID       string `json:"pid"`               //The upper ID
	SuperPID  string `json:"superID,omitempty"` //ID of the top-level parent
	Grade     int    `json:"grade,omitempty"`   //Department level
	Attr      int    `json:"attr"`              //1 company, 2 department
}

//GetUserInfo GetUserInfo
func (u *user) GetUserInfo(ctx context.Context, r *OneUserRequest) (*OneUserResponse, error) {
	response := &OneUserResponse{}
	err := client.POST(ctx, &u.client, host+oneUserURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

//GetUserByIDsRequest get user by ids request
type GetUserByIDsRequest struct {
	IDs []string `json:"ids"`
}

// GetUserByIDsResponse get user by ids response
type GetUserByIDsResponse struct {
	Users []OneUserResponse `json:"users"`
}

//GetUserByIDs Get user by ids
func (u *user) GetUserByIDs(ctx context.Context, r *GetUserByIDsRequest) (*GetUserByIDsResponse, error) {
	response := &GetUserByIDsResponse{}
	err := client.POST(ctx, &u.client, host+usersByIDsURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GetDepByIDsRequest GetDepByIDsRequest
type GetDepByIDsRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

// GetDepByIDsResponse GetDepByIDsResponse
type GetDepByIDsResponse struct {
	Deps []DepOneResponse `json:"deps"`
}

//GetDepByIDs Batch query departments by user IDs
func (u *user) GetDepByIDs(ctx context.Context, r *GetDepByIDsRequest) (*GetDepByIDsResponse, error) {
	response := &GetDepByIDsResponse{}
	err := client.POST(ctx, &u.client, host+depByIDsURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GetUsersByDepIDRequest GetUsersByDepIDRequest
type GetUsersByDepIDRequest struct {
	DepID          string `json:"depID"`
	IsIncludeChild int    `json:"isIncludeChild"`
}

// GetUsersByDepIDResponse GetUsersByDepIDResponse
type GetUsersByDepIDResponse struct {
	Users []OneUserResponse `json:"users"`
}

// GetUsersByDepID Query the users by department ID
func (u *user) GetUsersByDepID(ctx context.Context, r *GetUsersByDepIDRequest) (*GetUsersByDepIDResponse, error) {
	response := &GetUsersByDepIDResponse{}
	err := client.POST(ctx, &u.client, host+usersByDepIDURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}

// GetDepMaxGradeRequest request
type GetDepMaxGradeRequest struct {
}

// GetDepMaxGradeResponse response
type GetDepMaxGradeResponse struct {
	Grade int64 `json:"grade"`
}

//GetDepMaxGrade GetDepMaxGrade
func (u *user) GetDepMaxGrade(ctx context.Context, r *GetDepMaxGradeRequest) (*GetDepMaxGradeResponse, error) {
	response := &GetDepMaxGradeResponse{}
	err := client.POST(ctx, &u.client, host+depMaxGradeURI, r, response)
	if err != nil {
		return nil, err
	}
	return response, err
}
