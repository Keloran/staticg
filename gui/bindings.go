package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func BindKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return fmt.Errorf("failed to bind quit: %v", err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, save); err != nil {
		return fmt.Errorf("failed to bind save: %w", err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, preview); err != nil {
		return fmt.Errorf("failed to bind preview: %w", err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func save(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func preview(g *gocui.Gui, v *gocui.View) error {
	return nil
}
