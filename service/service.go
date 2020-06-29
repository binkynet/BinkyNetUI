//    Copyright 2017 Ewout Prangsma
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package service

import (
	"context"

	api "github.com/binkynet/BinkyNet/apis/v1"
	"github.com/binkynet/BinkyNet/discovery"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// Service is the API exposed by this service.
type Service interface {
	// Run the service until the given context is cancelled.
	Run(ctx context.Context) error
}

type Config struct {
	// LocalWorker version (semver) that is expected.
	// If the actual version is different, the LocalWorker must update
	// itself.
	RequiredWorkerVersion string
}

type Dependencies struct {
	Log zerolog.Logger
}

type service struct {
	log zerolog.Logger
	Config

	nwCtrlListener *discovery.ServiceListener
}

// NewService creates a Service instance and returns it.
func NewService(conf Config, deps Dependencies) (Service, error) {
	log := deps.Log.With().Str("component", "service").Logger()
	s := &service{
		log:    log,
		Config: conf,
	}
	s.nwCtrlListener = discovery.NewServiceListener(log, api.ServiceTypeNetworkControl, true, s.onNetworkControlChanged)
	return s, nil
}

// Run the manager until the given context is cancelled.
func (s *service) Run(ctx context.Context) error {
	log := s.log
	defer func() {
		log.Debug().Msg("Run finished")
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return s.nwCtrlListener.Run(ctx) })

	return g.Wait()
}

func (s *service) onNetworkControlChanged(info api.ServiceInfo) {

}
