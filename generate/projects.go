package generate

import (
	"fmt"
)

func CurrentProjects(root string) ([]File, error) {
	files, err := GetFiles(root + "projects/current")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title:    "Current Projects",
		NewIndex: root + "projects/current/newIndex.md",
		Index:    root + "projects/current/index.md",
		Path:     "/projects/current",
		Pages:    files,
	}
	err = p.Generate()
	if err != nil {
		return files, fmt.Errorf("generate template: %w", err)
	}

	if len(files) >= 5 {
		return files[:4], nil
	}

	return files, nil
}

func PastProjects(root string) ([]File, error) {
	files, err := GetFiles(root + "projects/past")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title:    "Past Projects",
		NewIndex: root + "projects/past/newIndex.md",
		Index:    root + "projects/past/index.md",
		Path:     "/projects/past",
		Pages:    files,
	}
	err = p.Generate()
	if err != nil {
		return files, fmt.Errorf("generate template: %w", err)
	}

	if len(files) >= 5 {
		return files[:4], nil
	}

	return files, nil
}
