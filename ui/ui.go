// Copyright 2020 Ewout Prangsma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author Ewout Prangsma
//

package ui

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/app"

	api "github.com/binkynet/BinkyNet/apis/v1"
	"github.com/rs/zerolog"
)

type UI struct {
	log        zerolog.Logger
	app        fyne.App
	mainWindow fyne.Window
	mainPage   *mainPage
}

// NewUI initialize a new UI
func NewUI(ctx context.Context, log zerolog.Logger) (*UI, error) {
	a := app.New()
	mainPage := NewMainPage(ctx, log)
	mainWindow := a.NewWindow("BinkyNet")
	mainWindow.SetContent(mainPage.Root)
	mainWindow.Resize(fyne.NewSize(800, 600))

	ui := &UI{
		log:        log,
		app:        a,
		mainWindow: mainWindow,
		mainPage:   mainPage,
	}
	return ui, nil
}

// Run until the given context is canceled
func (ui *UI) Run(ctx context.Context) error {
	ui.mainWindow.ShowAndRun()

	return nil
}

// CommandStationChanged is called when a new CommandStation service is detected.
func (ui *UI) CommandStationChanged(ctx context.Context, apic api.CommandStationServiceClient) {
	ui.mainPage.CommandStationChanged(ctx, apic)
}

// NetworkControlChanged is called when a new NetworkControl service is detected.
func (ui *UI) NetworkControlChanged(ctx context.Context, apic api.NetworkControlServiceClient) {
	ui.mainPage.NetworkControlChanged(ctx, apic)
}
