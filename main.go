package main

import (
	"fmt"
	"sync"

	"github.com/keloran/staticg/generate"
	"github.com/keloran/staticg/gui"
)

func main() {
  if err := gui.StartGUI(); err != nil {
    fmt.Printf("failed to start GUI: %v\n", err)
    return
  }

	// err = _main(os.Args[1:])
	// if err != nil {
	// 	fmt.Printf("failed: %v\n", err)
	// 	return
	// }
}

func _main(args []string) error {
	root := "./"
	if len(args) >= 1 {
		if args[0] != "" {
			root = args[0] + "/"
		}
	}

	errChan := make(chan error)
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)
	go func() {
		waitGroup.Wait()
		close(errChan)
	}()

	fmt.Printf("Create new item y/n ? ")
	create, err := getResponse(NEWLINE)
	if create[0] == YES {
		if err := editor(root); err != nil {
			errChan <- fmt.Errorf("editor err: %w", err)
		}
	}
	if err != nil {
		errChan <- fmt.Errorf("item question: %w", err)
	}

	ic := generate.IndexContent{}

	go func() {
		defer waitGroup.Done()
		ic, err = generate.Pages(root)
		if err != nil {
			errChan <- fmt.Errorf("generate error: %w", err)
		}
	}()

	go func() {
		defer waitGroup.Done()
		if err := ic.Generate(); err != nil {
			errChan <- fmt.Errorf("generate index: %w", err)
		}
		fmt.Printf("index categories generated\n")
	}()

	go func() {
		defer waitGroup.Done()
		if err := ic.GenerateFeed(); err != nil {
			errChan <- fmt.Errorf("generate feed: %w", err)
		}
		fmt.Printf("rss generated\n")
	}()

	waitGroup.Wait()
	for err := range errChan {
		if err != nil {
			return fmt.Errorf("main err: %w", err)
		}
	}

	return nil
}
