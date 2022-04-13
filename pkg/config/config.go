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

package config

import (
	"io/ioutil"
	"time"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"

	"github.com/quanxiang-cloud/cabin/logger"
	"gopkg.in/yaml.v2"
)

// Config Config
var Config *Configs

// DefaultPath DefaultPath
var DefaultPath = "./configs/config.yml"

// AppCenter AppCenter
type AppCenter struct {
	Model      string     `yaml:"model"`
	HTTPServer HTTPServer `yaml:"http"`
}

// Chaos Chaos
type Chaos struct {
	Model      string     `yaml:"model"`
	HTTPServer HTTPServer `yaml:"http"`
}

// Configs Configs
type Configs struct {
	AppCenter AppCenter `yaml:"app-center"`
	Chaos     Chaos     `yaml:"chaos"`

	Model      string
	HTTPServer HTTPServer

	Mysql             mysql2.Config   `yaml:"mysql"`
	Log               logger.Config   `yaml:"log"`
	InternalNet       client.Config   `yaml:"internalNet"`
	Redis             redis2.Config   `yaml:"redis"`
	InnerHost         InnerHostConfig `yaml:"innerHost"`
	CompatibleVersion string          `yaml:"compatibleVersion"`

	InitServerBits int `yaml:"initServerBits"`

	WorkLoad     int               `yaml:"workLoad"`
	MaximumRetry int               `yaml:"maximumRetry"`
	WaitTime     int               `yaml:"waitTime"`
	CachePath    string            `yaml:"cachePath"`
	KV           map[string]string `yaml:"kv"`
}

// InnerHostConfig InnerHostConfig
type InnerHostConfig struct {
	StructorHost string `yaml:"structor"`
	FlowHost     string `yaml:"flow"`
	PolyAPI      string `yaml:"polyAPI"`
	OrgHost      string `yaml:"org"`
}

// HTTPServer HTTPServer
type HTTPServer struct {
	Port              string        `yaml:"port"`
	ReadHeaderTimeOut time.Duration `yaml:"readHeaderTimeOut"`
	WriteTimeOut      time.Duration `yaml:"writeTimeOut"`
	MaxHeaderBytes    int           `yaml:"maxHeaderBytes"`
}

// Init Init
func Init(configPath string) error {
	if configPath == "" {
		configPath = "../configs/configs.yml"
	}
	Config = new(Configs)
	err := read(configPath, Config)
	if err != nil {
		return err
	}
	return nil
}

// NewConfig NewConfig
func NewConfig(path string) (*Configs, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		return nil, err
	}

	return Config, nil
}

// read read
func read(yamlPath string, v interface{}) error {
	// Read config file
	buf, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	return nil
}
