/*
 * Created on Mon Jun 29 2020
 *
 * Copyright (c) 2020 Your Company
 */
package ui

import (
	"context"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type UI struct {
	theme *material.Theme
}

// NewUI initialize a new UI
func NewUI() (*UI, error) {
	ui := &UI{
		theme: material.NewTheme(gofont.Collection()),
	}
	return ui, nil
}

// Run until the given context is canceled
func (ui *UI) Run(ctx context.Context) error {
	w := app.NewWindow()
	var ops op.Ops

	for {
		select {
		case <-ctx.Done():
			// Context canceled
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				ui.layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}

}

func (ui *UI) layout(gtx layout.Context) {
	material.H1(ui.theme, "Hello Binky")
}
