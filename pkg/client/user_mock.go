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

	"github.com/stretchr/testify/mock"
)

type userMock struct {
	mock.Mock
}

// NewUserMock NewUserMock
func NewUserMock() User {
	// create an instance of our test object
	user := new(userMock)

	// setup expectations
	user.On("GetInfo", mock.Anything).Return(map[string]UserInfo{
		"-1": {
			ID:       "-1",
			UserName: "alex",
		},
	}, nil)

	user.On("GetDepartment", mock.Anything).Return(map[string]Department{
		"-1": {
			ID:             "-1",
			DepartmentName: "department",
		},
	}, nil)

	user.On("GetUsersByDEPID", mock.Anything).Return(map[string][]UserInfo{
		"-1": {
			{
				ID:          "-1",
				UserName:    "alex",
				Email:       "test@test.com",
				Phone:       "13988886666",
				Avatar:      "xxxxxx",
				IsDEPLeader: 1,
			}},
	}, nil)

	return user
}

func (m *userMock) GetInfo(ctx context.Context, userIDs ...string) ([]UserInfo, error) {
	args := m.Called()
	userMap := args.Get(0).(map[string]UserInfo)
	ans := make([]UserInfo, 0, len(userIDs))
	for _, userID := range userIDs {
		userInfo, ok := userMap[userID]
		if !ok {
			// Return null data
			ans = append(ans, UserInfo{})
			continue
		}
		ans = append(ans, userInfo)
	}

	return ans, args.Error(1)
}

func (m *userMock) GetDepartment(ctx context.Context, ids ...string) ([]Department, error) {
	args := m.Called()
	departmentMap := args.Get(0).(map[string]Department)
	ans := make([]Department, 0, len(ids))
	for _, userID := range ids {
		department, ok := departmentMap[userID]
		if !ok {
			// Return null data
			ans = append(ans, Department{})
			continue
		}
		ans = append(ans, department)
	}

	return ans, args.Error(1)
}

func (m *userMock) GetUsersByDEPID(ctx context.Context, depID string, includeChildDEPChild, page, limit int) ([]UserInfo, error) {
	args := m.Called()
	userMap := args.Get(0).(map[string]UserInfo)
	ans := make([]UserInfo, 0)
	userInfo, ok := userMap[depID]
	if !ok {
		// Return null data
		ans = append(ans, UserInfo{})

	}
	ans = append(ans, userInfo)

	return ans, args.Error(1)
}
