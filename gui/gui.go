package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func StartGUI() error {
  g, err := gocui.NewGui(gocui.OutputNormal)
  if err != nil {
    return fmt.Errorf("failed to create window: %w", err)
  }
  defer g.Close()

  g.Cursor = true
  g.Mouse = true

  g.SetManagerFunc(layout)
  if err := BindKeys(g); err != nil {
    return fmt.Errorf("failed to bind keys: %w", err)
  }

  if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
    return fmt.Errorf("failed to start main loop: %w", err)
  }

  return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if _, err := g.SetView("side", -1, -1, int(0.2 * float32(maxX)), maxY - 5); err != nil && err != gocui.ErrUnknownView {
		return fmt.Errorf("failed to create side: %w", err)
	}

	if _, err := g.SetView("editor", int(0.2 * float32(maxX)), -1, maxX, maxY - 5); err != nil && err != gocui.ErrUnknownView {
		return fmt.Errorf("failed to create editor: %w", err)
	}

	if _, err := g.SetView("cmd", -1, maxY-5, maxX, maxX); err != nil && err != gocui.ErrUnknownView {
		return fmt.Errorf("failed to create cmd: %w", err)
	}

	// if v, err := g.SetView("but1", 2, 2, 22, 7); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return fmt.Errorf("failed to set view but1: %w", err)
	// 	}
	// 	v.Highlight = true
	// 	v.SelBgColor = gocui.ColorGreen
	// 	v.SelFgColor = gocui.ColorBlack
	// 	_, _ = fmt.Fprintln(v, "Button 1 - line 1")
	// 	_, _ = fmt.Fprintln(v, "Button 1 - line 2")
	// 	_, _ = fmt.Fprintln(v, "Button 1 - line 3")
	// 	_, _ = fmt.Fprintln(v, "Button 1 - line 4")
	// }
  //
	// if v, err := g.SetView("but2", 24, 2, 44, 4); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return fmt.Errorf("failed to set but2 view: %w", err)
	// 	}
	// 	v.Highlight = true
	// 	v.SelBgColor = gocui.ColorGreen
	// 	v.SelFgColor = gocui.ColorBlack
	// 	_, _ = fmt.Fprintln(v, "Button 2 - line 1")
	// }

	return nil
}
