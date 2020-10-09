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
	"fyne.io/fyne/widget"

	api "github.com/binkynet/BinkyNet/apis/v1"
	"github.com/rs/zerolog"
)

type mainPage struct {
	widget.Box
}

// NewMainPage constructs a new main UI page.
func NewMainPage(ctx context.Context, log zerolog.Logger, apic api.NetworkControlServiceClient) fyne.CanvasObject {
	powerPanel, powerItems := NewPowerPanel(ctx, log, apic)
	toolbar := widget.NewToolbar(
		powerItems...,
	)
	return &mainPage{
		Box: *widget.NewVBox(toolbar, powerPanel),
	}
}
