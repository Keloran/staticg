package main

import (
	"fmt"
)

func main() {
	err := blogPages()
	if err != nil {
		fmt.Printf("blog err: %v\n", err)
		return
	}

	err = currentProjects()
	if err != nil {
		fmt.Printf("current projects: %v\n", err)
		return
	}

	err = pastProjects()
	if err != nil {
		fmt.Printf("past projects: %v\n", err)
		return
	}
}

func blogPages() error {
	return nil
}

func currentProjects() error {
	return nil
}

func pastProjects() error {
	return nil
}

