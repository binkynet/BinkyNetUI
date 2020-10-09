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
	"net"
	"strconv"

	api "github.com/binkynet/BinkyNet/apis/v1"
	"github.com/binkynet/BinkyNet/discovery"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Service is the API exposed by this service.
type Service interface {
	// Run the service until the given context is cancelled.
	Run(ctx context.Context) error
}

// DiscoveryListener is informed when certain discovery events happen.
type DiscoveryListener interface {
	NetworkControlChanged(ctx context.Context, api api.NetworkControlServiceClient)
}

type Config struct {
	// LocalWorker version (semver) that is expected.
	// If the actual version is different, the LocalWorker must update
	// itself.
	RequiredWorkerVersion string
}

type Dependencies struct {
	Log               zerolog.Logger
	DiscoveryListener DiscoveryListener
}

type service struct {
	log               zerolog.Logger
	discoveryListener DiscoveryListener
	Config

	nwCtrlListener *discovery.ServiceListener
	nwCtrlChanges  chan api.ServiceInfo
}

// NewService creates a Service instance and returns it.
func NewService(conf Config, deps Dependencies) (Service, error) {
	log := deps.Log.With().Str("component", "service").Logger()
	s := &service{
		log:               log,
		discoveryListener: deps.DiscoveryListener,
		Config:            conf,
	}
	s.nwCtrlChanges = make(chan api.ServiceInfo, 8)
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
	g.Go(func() error { return s.run(ctx) })

	return g.Wait()
}

// run the actual service
func (s *service) run(ctx context.Context) error {
	var nwCtrlConn grpcConn
	for {
		select {
		case <-ctx.Done():
			// Context canceled
			nwCtrlConn.Close()
			return nil
		case info := <-s.nwCtrlChanges:
			log := s.log.With().
				Str("address", info.GetApiAddress()).
				Logger()
			log.Debug().Msg("NetworkControl service changed")
			conn, err := dialConn(&info)
			if err != nil {
				log.Warn().Err(err).Msg("Dialing NetworkControl failed")
				continue
			}
			nwCtrlConn.Close()
			nwCtrlConn.SetConn(ctx, conn)
			nwCtrlAPI := api.NewNetworkControlServiceClient(conn)
			s.discoveryListener.NetworkControlChanged(nwCtrlConn.ctx, nwCtrlAPI)
		}
	}
}

// NetworkControl service has changed
func (s *service) onNetworkControlChanged(info api.ServiceInfo) {
	s.nwCtrlChanges <- info
}

func dialConn(info *api.ServiceInfo) (*grpc.ClientConn, error) {
	address := net.JoinHostPort(info.GetApiAddress(), strconv.Itoa(int(info.GetApiPort())))
	var opts []grpc.DialOption
	if !info.Secure {
		opts = append(opts, grpc.WithInsecure())
	}
	return grpc.Dial(address, opts...)
}
