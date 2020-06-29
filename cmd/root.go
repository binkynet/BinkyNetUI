//    Copyright 2020 Ewout Prangsma
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

package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"gioui.org/app"

	"github.com/binkynet/BinkyNetUI/service"
	"github.com/binkynet/BinkyNetUI/ui"
)

const (
	projectName     = "BinkyNet UI"
	defaultGrpcPort = 8823
)

var (
	// RootCmd is the root command of the program
	RootCmd = &cobra.Command{
		Short: "BinkyNET UI",
		Run:   runRootCmd,
	}
	rootArgs struct {
		service        service.Config
		projectVersion string
		projectBuild   string
	}
	cliLog = zerolog.New(os.Stdout)
)

// SetVersion records given version info
func SetVersion(version, build string) {
	rootArgs.projectVersion = version
	rootArgs.projectBuild = build
}

func init() {
}

func runRootCmd(cmd *cobra.Command, args []string) {
	svc, err := service.NewService(rootArgs.service, service.Dependencies{
		Log: cliLog,
	})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("NewService failed")
	}
	ui, err := ui.NewUI()
	if err != nil {
		cliLog.Fatal().Err(err).Msg("NewUI failed")
	}
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return svc.Run(ctx) })
	g.Go(func() error { return ui.Run(ctx) })
	go func() {
	if err := g.Wait(); err != nil {
		cliLog.Fatal().Err(err).Msg("Run failed")
	}
}()

	app.Main()
}
