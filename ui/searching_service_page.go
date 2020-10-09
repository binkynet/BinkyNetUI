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
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)

func NewSearchingServicePage() fyne.CanvasObject {
	circle := canvas.NewCircle(color.White)
	circle.StrokeColor = color.Gray{0x99}
	circle.StrokeWidth = 5
	circle.Resize(fyne.Size{200, 300})

	lb := canvas.NewText("Searching for services...", color.RGBA{255, 0, 0, 128})

	c := fyne.NewContainerWithLayout(layout.NewMaxLayout(), circle, fyne.NewContainerWithLayout(layout.NewCenterLayout(), lb))

	return c
}
