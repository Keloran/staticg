package main

import (
	"fmt"
)

func currentProjects(root string) ([]File, error) {
	files, err := getFiles(root + "projects/current")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title:    "Current Projects",
		NewIndex: root + "projects/current/newIndex.md",
		Index:    root + "projects/current/index.md",
		Pages:    files,
	}
	err = p.generate()
	if err != nil {
		return files, fmt.Errorf("generate template: %w", err)
	}

	if len(files) >= 5 {
		return files[:4], nil
	}

	return files, nil
}

func pastProjects(root string) ([]File, error) {
	files, err := getFiles(root + "projects/past")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title:    "Past Projects",
		NewIndex: root + "projects/past/newIndex.md",
		Index:    root + "projects/past/index.md",
		Pages:    files,
	}
	err = p.generate()
	if err != nil {
		return files, fmt.Errorf("generate template: %w", err)
	}

	if len(files) >= 5 {
		return files[:4], nil
	}

	return files, nil
}
