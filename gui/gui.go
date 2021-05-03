package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func CreateGUI() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return g, fmt.Errorf("failed to create window: %w", err)
	}

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)
	if err := BindKeys(g); err != nil {
		return g, fmt.Errorf("failed to bind keys: %w", err)
	}

	return g, nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("help", -1, -1, 50, maxY - 5); err != nil {
	  if err != gocui.ErrUnknownView {
	    return fmt.Errorf("failed to create help view: %w", err)
    }
    if err := drawHelp(v); err != nil {
      return fmt.Errorf("failed to update help view: %w", err)
    }
  }

	if v, err := g.SetView("commands", -1, maxY - 5, 50, maxY); err != nil {
	  if err != gocui.ErrUnknownView {
      return fmt.Errorf("failed to create commands view: %w", err)
    }

    if err := drawCommands(v); err != nil {
      return fmt.Errorf("failed to update commands view: %w", err)
    }
	}

	if v, err := g.SetView("editor", 50, -1, maxX, maxY); err != nil {
	  if err != gocui.ErrUnknownView {
      return fmt.Errorf("failed to create editor: %w", err)
    }

    v.Editable = true
    v.Wrap = true
    if _, err := g.SetCurrentView("editor"); err != nil {
      return fmt.Errorf("failed to set editor as current view: %w", err)
    }
  }

	return nil
}

func drawCommands(v *gocui.View) error {
	if _, err := fmt.Fprintln(v, "Quit: ^q"); err != nil {
	  return fmt.Errorf("failed to write quit to commands buffer: %w", err)
  }
	if _, err := fmt.Fprintln(v, "Preview: ^p"); err != nil {
    return fmt.Errorf("failed to write preview to commands  buffer: %w", err)
  }
	if _, err := fmt.Fprintln(v, "Save: ^s"); err != nil {
    return fmt.Errorf("failed to write save to commands buffer: %w", err)
  }

	return nil
}

func drawHelp(v *gocui.View) error {
  if _, err := fmt.Fprintf(v, "# Heading\n\n"); err != nil {
    return fmt.Errorf("failed to write headings to help buffer: %w", err)
  }
  if _, err := fmt.Fprintf(v, "* List\n\n"); err != nil {
    return fmt.Errorf("failed to write list to help buffer: %w", err)
  }
  if _, err := fmt.Fprintf(v, "[text](link) Link\n\n"); err != nil {
    return fmt.Errorf("failed to write links to help buffer: %w", err)
  }
  if _, err := fmt.Fprintf(v, "![alt](url) Image\n\n"); err != nil {
    return fmt.Errorf("failed to write image to help buffer: %w", err)
  }
  if _, err := fmt.Fprintf(v, "```<language>\ncode```\n\n"); err != nil {
    return fmt.Errorf("failed to write code to help buffer: %w", err)
  }
  if _, err := fmt.Fprintf(v, "--- Line\n\n"); err != nil {
    return fmt.Errorf("failed to write line to help buffer: %w", err)
  }
  if _, err := fmt.Fprintf(v, "[![alt](image)](video link) Video\n\n"); err != nil {
    return fmt.Errorf("failed to write video to help buffer: %w", err)
  }

  return nil
}

func StartGUI(g *gocui.Gui) error {
  defer g.Close()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("failed to start main loop: %w", err)
	}

	return nil
}
