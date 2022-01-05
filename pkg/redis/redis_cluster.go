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

package redis

import (
	"github.com/quanxiang-cloud/appcenter/pkg/config"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"

	"github.com/go-redis/redis/v8"
)

// ClusterClient redis cluster client
var ClusterClient *redis.ClusterClient

// Init Init
func Init() error {
	return newClient()
}

func newClient() error {
	client, err := redis2.NewClient(config.Config.Redis)
	if err != nil {
		return err
	}
	ClusterClient = client
	return nil
}

// Close Close
func Close() error {
	return ClusterClient.Close()
}
