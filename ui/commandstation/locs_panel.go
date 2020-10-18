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

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	api "github.com/binkynet/BinkyNet/apis/v1"
)

type locsPanel struct {
	locList       *widget.Group
	selectedPanel *fyne.Container

	requests chan *api.Loc
}

func NewLocsPanel(ctx context.Context, log zerolog.Logger, apic api.CommandStationServiceClient) (fyne.CanvasObject, []widget.ToolbarItem) {
	p := &locsPanel{
		locList:       widget.NewGroupWithScroller("Loc list", widget.NewLabel("TODO")),
		selectedPanel: fyne.NewContainer(),
		requests:      make(chan *api.Loc, 8),
	}
	p.updateState(false, false, true)

	splitContainer := widget.NewHSplitContainer(widget.NewVScrollContainer(p.locList), p.selectedPanel)
	splitContainer.SetOffset(0.3)

	go p.run(ctx, log, apic)

	return splitContainer, []widget.ToolbarItem{}
}

func (p *locsPanel) updateState(requested, actual, unknown bool) {
	//util.SetVisible(p.tbbPowerOff, unknown || actual || changing)
	//util.SetVisible(p.tbbPowerOn, unknown || !actual || changing)
}

func (p *locsPanel) run(ctx context.Context, log zerolog.Logger, apic api.CommandStationServiceClient) {
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

func (p *locsPanel) runOnce(ctx context.Context, log zerolog.Logger, apic api.CommandStationServiceClient) error {
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
