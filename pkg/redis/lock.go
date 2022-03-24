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
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Locker locker
type Locker struct {
	Key        string
	Value      string
	ExpireTime time.Duration
	Conn       redis.UniversalClient
	Ctx        context.Context
}

// NewLocker new
func NewLocker(key, value string, expireTime time.Duration, conn redis.UniversalClient) *Locker {
	return &Locker{Key: key, Value: value, ExpireTime: expireTime, Conn: conn, Ctx: context.Background()}
}

// Lock lock action
func (o *Locker) Lock() (bool, error) {

	_, err := o.Conn.SetNX(o.Ctx, o.Key, o.Value, o.ExpireTime*time.Second).Result()
	if err != nil {
		return false, err
	}
	get, err := o.Conn.Get(o.Ctx, o.Key).Result()
	if err != nil {
		return false, err
	}
	if get != "" && get == o.Value {
		return true, nil
	}
	return false, nil

}

// UnLock unlock
func (o *Locker) UnLock() error {
	_, err := o.Conn.Del(o.Ctx, o.Key).Result()
	return err
}
