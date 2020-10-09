/*
 * Created on Mon Jun 29 2020
 *
 * Copyright (c) 2020 Your Company
 */
package ui

import (
	"context"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
)

type UI struct {
}

// NewUI initialize a new UI
func NewUI() (*UI, error) {
	ui := &UI{}
	return ui, nil
}

// Run until the given context is canceled
func (ui *UI) Run(ctx context.Context) error {
	myApp := app.New()
	myWindow := myApp.NewWindow("Form Layout")

	label1 := canvas.NewText("Label 1", color.Black)
	value1 := canvas.NewText("Value", color.White)
	label2 := canvas.NewText("Label 2", color.Black)
	value2 := canvas.NewText("Something", color.White)
	grid := fyne.NewContainerWithLayout(layout.NewFormLayout(),
		label1, value1, label2, value2)
	myWindow.SetContent(grid)
	myWindow.ShowAndRun()

	return nil
}
