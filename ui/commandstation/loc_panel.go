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

package commandstation

import (
	"context"
	"image/color"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	api "github.com/binkynet/BinkyNet/apis/v1"
)

type locPanel struct {
	widget.Box

	lbName    *canvas.Text
	tbAddress *widget.Entry
	slSpeed   *widget.Slider

	requests chan *api.Loc

	address    int
	speed      int
	speedSteps int32
	direction  api.LocDirection
	light      bool
}

func NewLocPanel(ctx context.Context, log zerolog.Logger, apic api.CommandStationServiceClient) fyne.CanvasObject {
	p := &locPanel{
		lbName:     canvas.NewText("Name", color.RGBA{0, 0, 255, 255}),
		tbAddress:  widget.NewEntry(),
		slSpeed:    widget.NewSlider(0, 100),
		requests:   make(chan *api.Loc, 8),
		address:    3,
		speed:      0,
		speedSteps: 128,
		direction:  api.LocDirection_FORWARD,
		light:      true,
	}
	p.tbAddress.OnChanged = p.onAddressChanged
	p.tbAddress.SetText("3")
	p.slSpeed.OnChanged = p.onSpeedChanged
	chDirection := widget.NewRadio([]string{"Forward", "Reverse"}, func(x string) {
		if x == "Forward" {
			p.onDirectionChanged(api.LocDirection_FORWARD)
		} else {
			p.onDirectionChanged(api.LocDirection_REVERSE)
		}
	})
	chDirection.SetSelected("Forward")
	cbLight := widget.NewCheck("Light", func(value bool) {
		p.light = value
		p.sendRequest()
	})
	p.Box = *widget.NewVBox(p.lbName, p.tbAddress, p.slSpeed, chDirection, cbLight)

	go p.run(ctx, log, apic)

	return p
}

func (p *locPanel) sendRequest() {
	p.requests <- &api.Loc{
		Address: api.ObjectAddress(strconv.Itoa(p.address)),
		Request: &api.LocState{
			Speed:      int32(p.speed),
			SpeedSteps: p.speedSteps,
			Direction:  p.direction,
			Functions: map[int32]bool{
				0: p.light,
			},
		},
	}
}

func (p *locPanel) onAddressChanged(value string) {
	if address, err := strconv.Atoi(value); err == nil {
		p.address = address
		p.sendRequest()
	} else {
		p.tbAddress.SetText(strconv.Itoa(p.address))
	}
}

func (p *locPanel) onSpeedChanged(value float64) {
	p.speed = int(value)
	p.sendRequest()
}

func (p *locPanel) onDirectionChanged(value api.LocDirection) {
	p.direction = value
	p.sendRequest()
}

func (p *locPanel) run(ctx context.Context, log zerolog.Logger, apic api.CommandStationServiceClient) {
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

func (p *locPanel) runOnce(ctx context.Context, log zerolog.Logger, apic api.CommandStationServiceClient) error {
	locStream, err := apic.Locs(ctx)
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	// Receive loop
	g.Go(func() error {
		for {
			msg, err := locStream.Recv()
			if err != nil {
				log.Debug().Err(err).Msg("locStream.Recv failed")
				return err
			}
			log.Debug().Interface("msg", msg).Msg("Received loc event")
			//			p.updateState(msg.GetRequest().GetEnabled(), msg.GetActual().GetEnabled(), false)
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
				if err := locStream.Send(req); err != nil {
					log.Debug().Err(err).Msg("locStream.Send failed")
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
