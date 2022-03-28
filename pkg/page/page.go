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

package page

const (
	defaultSize = 20
	defaultPage = 1
)

// Page Page
type Page struct {
	PageSize    int         `json:"-"`
	TotalCount  int64       `json:"total_count"`
	TotalPage   int         `json:"-"`
	CurrentPage int         `json:"-"`
	StartIndex  int         `json:"-"`
	Data        interface{} `json:"data"`
}

// NewPage NewPage
func NewPage(currentPage int, pageSize int, totalCount int64) *Page {
	page := Page{
		PageSize:    pageSize,
		CurrentPage: currentPage,
	}
	if pageSize == 0 {
		page.PageSize = defaultSize
	}
	if currentPage == 0 {
		page.CurrentPage = defaultPage
	}
	page.StartIndex = (page.CurrentPage - 1) * page.PageSize
	page.TotalCount = totalCount
	if page.TotalCount%int64(page.PageSize) == 0 {
		page.TotalPage = int(page.TotalCount) / page.PageSize
	} else {
		page.TotalPage = int(page.TotalCount)/page.PageSize + 1
	}
	return &page
}
