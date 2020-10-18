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

package networkcontrol

import (
	"context"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	api "github.com/binkynet/BinkyNet/apis/v1"
	"github.com/binkynet/BinkyNetUI/ui/util"
)

type powerPanel struct {
	lbPowerIsUnknown *canvas.Text
	lbPowerIsOff     *canvas.Text
	lbPowerIsOn      *canvas.Text
	lbPowerGoingOff  *canvas.Text
	lbPowerGoingOn   *canvas.Text
	butPowerOn       *widget.Button
	butPowerOff      *widget.Button

	tbbPowerOn  widget.ToolbarItem
	tbbPowerOff widget.ToolbarItem

	requests chan bool
}

func NewPowerPanel(ctx context.Context, log zerolog.Logger, apic api.NetworkControlServiceClient) (fyne.CanvasObject, []widget.ToolbarItem) {
	p := &powerPanel{
		lbPowerIsUnknown: canvas.NewText("Power is not known", color.RGBA{0, 0, 255, 255}),
		lbPowerIsOff:     canvas.NewText("Power is off", color.RGBA{255, 0, 0, 255}),
		lbPowerIsOn:      canvas.NewText("Power is on", color.RGBA{0, 255, 0, 255}),
		lbPowerGoingOff:  canvas.NewText("Power turning off...", color.RGBA{255, 0, 0, 255}),
		lbPowerGoingOn:   canvas.NewText("Power turning on...", color.RGBA{0, 255, 0, 255}),
		requests:         make(chan bool, 8),
	}
	p.butPowerOn = widget.NewButton("Power on!", func() {
		p.requests <- true
	})
	p.butPowerOff = widget.NewButton("Power off!", func() {
		p.requests <- false
	})
	p.tbbPowerOn = widget.NewToolbarAction(theme.VolumeUpIcon(), func() {
		p.requests <- true
	})
	p.tbbPowerOff = widget.NewToolbarAction(theme.VolumeMuteIcon(), func() {
		p.requests <- false
	})
	p.updateState(false, false, true)

	box := widget.NewVBox(p.lbPowerIsUnknown, p.lbPowerIsOff, p.lbPowerIsOn, p.lbPowerGoingOff, p.lbPowerGoingOn, p.butPowerOff, p.butPowerOn)

	go p.run(ctx, log, apic)

	return box, []widget.ToolbarItem{p.tbbPowerOn, p.tbbPowerOff}
}

func (p *powerPanel) updateState(requested, actual, unknown bool) {
	changing := requested != actual
	util.SetVisible(p.lbPowerIsUnknown, unknown)
	util.SetVisible(p.lbPowerIsOn, !unknown && !changing && (requested && actual))
	util.SetVisible(p.lbPowerIsOff, !unknown && !changing && (!requested && !actual))
	util.SetVisible(p.lbPowerGoingOn, !unknown && (requested && changing))
	util.SetVisible(p.lbPowerGoingOff, !unknown && (!requested && changing))
	util.SetVisible(p.butPowerOff, unknown || actual || changing)
	util.SetVisible(p.butPowerOn, unknown || !actual || changing)
	//util.SetVisible(p.tbbPowerOff, unknown || actual || changing)
	//util.SetVisible(p.tbbPowerOn, unknown || !actual || changing)
}

func (p *powerPanel) run(ctx context.Context, log zerolog.Logger, apic api.NetworkControlServiceClient) {
	defer close(p.requests)
	for {
		if ctx.Err() != nil {
			return
		}
		if err := p.runOnce(ctx, log, apic); err != nil {
			log.Warn().Err(err).Msg("runOnce failed")
		}
	}
}

func (p *powerPanel) runOnce(ctx context.Context, log zerolog.Logger, apic api.NetworkControlServiceClient) error {
	pwStream, err := apic.Power(ctx)
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	// Receive loop
	g.Go(func() error {
		for {
			msg, err := pwStream.Recv()
			if err != nil {
				log.Debug().Err(err).Msg("pwStream.Recv failed")
				return err
			}
			log.Debug().Msg("Received power event")
			p.updateState(msg.GetRequest().GetEnabled(), msg.GetActual().GetEnabled(), false)
		}
	})
	// Send loop
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				// Context canceled
				return nil
			case req := <-p.requests:
				// Send request
				if err := pwStream.Send(&api.PowerState{Enabled: req}); err != nil {
					log.Debug().Err(err).Msg("pwStream.Send failed")
					return err
				}
			}
		}
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
