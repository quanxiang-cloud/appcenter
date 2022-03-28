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

package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// InvalidURI Invalid URI
	InvalidURI = 90014000000
	// InvalidParams Invalid params
	InvalidParams = 90014000001
	// InvalidTimestamp Invalid timestamp
	InvalidTimestamp = 90014000002
	// NameExist Name already exist
	NameExist = 90014000003
	// InvalidDel Invalid Delection
	InvalidDel = 90014000004
	// ErrIdentifiesExist Identifies already exist
	ErrIdentifiesExist = 90014000005
	// ErrVersion Version incompatibility
	ErrVersion = 90014000006

	// ErrDataNotExist Data not exist
	ErrDataNotExist = 90014000007
	// ErrNoPermission Do not hava permission
	ErrNoPermission = 90014000008
	// ErrActionTimeOut Timeout
	ErrActionTimeOut = 90014000010
)

// CodeTable 码表
var CodeTable = map[int64]string{
	InvalidURI:         "无效的URI.",
	InvalidParams:      "无效的参数.",
	InvalidTimestamp:   "无效的时间格式.",
	NameExist:          "名称已被使用！请检查后重试！",
	InvalidDel:         "删除无效！对象不存在或请检查参数！",
	ErrIdentifiesExist: "唯一标识已存在",
	ErrVersion:         "版本不兼容",
	ErrDataNotExist:    "数据不存在",
	ErrNoPermission:    "没有权限",
	ErrActionTimeOut:   "操作超时，稍后请再次尝试",
}
