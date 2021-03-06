package generate

import (
	"fmt"
)

func BlogPages(root string) ([]File, error) {
	files, err := GetFiles(root + "blog")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title:    "Blog",
		NewIndex: root + "blog/newIndex.md",
		Index:    root + "blog/index.md",
		Path:     "/blog",
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
