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

package probe

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/quanxiang-cloud/cabin/logger"
)

const (
	readinessPending int32 = iota
	readinessTrue
	readinessFalse
)

// Probe probe
type Probe struct {
	readiness int32
}

// New return *Probe
func New() *Probe {
	return &Probe{
		// log:       log,
		readiness: readinessPending,
	}
}

func (p *Probe) setTrue() {
	atomic.StoreInt32(&p.readiness, readinessTrue)
}

func (p *Probe) setFalse() {
	atomic.StoreInt32(&p.readiness, readinessFalse)
}

func (p *Probe) getReadiness() int32 {
	return atomic.LoadInt32(&p.readiness)
}

// SetRunning set running
func (p *Probe) SetRunning() {
	logger.Logger.Info("probe ready")
	p.setTrue()
}

// LivenessProbe liveness probe
func (p *Probe) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	if p.getReadiness() != readinessFalse {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (p *Probe) isSafe(r *http.Request) bool {
	if strings.HasPrefix(r.Host, "127.0.0.1") ||
		strings.HasPrefix(r.Host, "localhost") {
		return true
	}
	return false
}

// ReadinessProbe readiness probe
func (p *Probe) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-readiness-shutdown") != "" {
		if !p.isSafe(r) {
			logger.Logger.Info("try to shutdown,but is not safe. refuse!", "host", r.Host)
			return
		}
		logger.Logger.Info("readiness shutdown")
		p.setFalse()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if p.getReadiness() == readinessTrue {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
