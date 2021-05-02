package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func BindKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, quit); err != nil {
		return fmt.Errorf("failed to bind quit: %v", err)
	}

	for _, n := range []string{"but1", "but2"} {
		if err := g.SetKeybinding(n, gocui.MouseLeft, gocui.ModNone, showMsg); err != nil {
			return fmt.Errorf("failed to bind showMsg mouse buttons: %w", err)
		}
	}

	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg); err != nil {
		return fmt.Errorf("failed to bind delMsg mouse: %w", err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return fmt.Errorf("failed to set current view: %w", err)
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", ((maxX / 2) - 10), maxY/2, ((maxX / 2) + 10), ((maxY / 2) + 2)); err != nil {
		if err != gocui.ErrUnknownView {
			return fmt.Errorf("failed to set msg view: %w", err)
		}
		if _, err := fmt.Fprintln(v, l); err != nil {
			return fmt.Errorf("failed to write to msg view: %w", err)
		}
	}

	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return fmt.Errorf("failed to delete msg view: %w", err)
	}

	return nil
}
