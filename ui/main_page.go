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

	"fyne.io/fyne/widget"
	"github.com/rs/zerolog"

	api "github.com/binkynet/BinkyNet/apis/v1"
	"github.com/binkynet/BinkyNetUI/ui/commandstation"
	"github.com/binkynet/BinkyNetUI/ui/networkcontrol"
)

type mainPage struct {
	log  zerolog.Logger
	Root *widget.TabContainer
	tabs []*widget.TabItem
}

// NewMainPage constructs a new main UI page.
func NewMainPage(ctx context.Context, log zerolog.Logger) *mainPage {
	tabs := []*widget.TabItem{
		widget.NewTabItem("Network control", NewSearchingServicePage("network control")),
		widget.NewTabItem("Command station", NewSearchingServicePage("command station")),
	}
	tc := widget.NewTabContainer(tabs...)
	return &mainPage{
		log:  log,
		Root: tc,
		tabs: tabs,
	}
}

// CommandStationChanged is called when a new CommandStation service is detected.
func (m *mainPage) CommandStationChanged(ctx context.Context, apic api.CommandStationServiceClient) {
	m.tabs[1].Content = commandstation.NewMainPage(ctx, m.log, apic)
	m.Root.Refresh()
	m.Root.SelectTabIndex(1)
}

// NetworkControlChanged is called when a new NetworkControl service is detected.
func (m *mainPage) NetworkControlChanged(ctx context.Context, apic api.NetworkControlServiceClient) {
	m.tabs[0].Content = networkcontrol.NewMainPage(ctx, m.log, apic)
	m.Root.Refresh()
	m.Root.SelectTabIndex(0)
}
