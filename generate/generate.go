package generate

import (
	"fmt"
	"os"
	"strings"
)

// IndexContent ...
type IndexContent struct {
	Blog    []File
	Current []File
	Past    []File
	Other   []File
}

// PageContent ...
type PageContent struct {
	Title    string
	NewIndex string
	Index    string
	Path     string
	Pages    []File
}

func writeContent(p *os.File, f []File, title, link string) error {
	if len(f) >= 1 {
		if link != "" {
			_, err := p.WriteString(fmt.Sprintf("### [%s](%s)\n", title, link))
			if err != nil {
				return fmt.Errorf("writeContent: %s title with link: %w", title, err)
			}
		} else {
			_, err := p.WriteString(fmt.Sprintf("### [%s]\n", title))
			if err != nil {
				return fmt.Errorf("writeContent %s title without link: %w", title, err)
			}
		}
		for _, f := range f {
			_, err := p.WriteString(fmt.Sprintf("* [%s](%s%s)\n", f.Title, link, f.CleanPath))
			if err != nil {
				return fmt.Errorf("writeContent: %s item: %+v, %w", title, f, err)
			}
		}
		_, err := p.WriteString("\n")
		if err != nil {
			return fmt.Errorf("writeContent closer: %w", err)
		}
	}

	return nil
}

func (ic IndexContent) Generate() error {
	r, err := os.Create("newIndex.md")
	if err != nil {
		return fmt.Errorf("create newindex: %w", err)
	}

	// Blog
	if fileExists("blog/index.md") {
		err = writeContent(r, ic.Blog, "Blog", "/blog")
		if err != nil {
			return fmt.Errorf("blog: %w", err)
		}
	}

	// Current
	if fileExists("projects/current/index.md") {
		err = writeContent(r, ic.Current, "Current Projects", "/projects/current")
		if err != nil {
			return fmt.Errorf("current projects: %w", err)
		}
	}

	// Past
	if fileExists("projects/past/index.md") {
		err = writeContent(r, ic.Past, "Past Projects", "/projects/past")
		if err != nil {
			return fmt.Errorf("past projects: %w", err)
		}
	}

	if len(ic.Other) >= 1 {
		err = writeContent(r, ic.Other, "Other Links", "")
		if err != nil {
			return fmt.Errorf("other links: %w", err)
		}
	}

	// Close file
	err = r.Close()
	if err != nil {
		return fmt.Errorf("index close: %w", err)
	}

	err = os.Rename("newIndex.md", "index.md")
	if err != nil {
		return fmt.Errorf("move index: %w", err)
	}

	return nil
}

func (p PageContent) Generate() error {
	r, err := os.Create(p.NewIndex)
	if err != nil {
		return fmt.Errorf("create index: %v, %w", p.NewIndex, err)
	}

	_, err = r.WriteString(fmt.Sprintf("### %s\n", p.Title))
	if err != nil {
		return fmt.Errorf("write title: %w", err)
	}

	for i := 1; i < len(p.Pages); i++ {
		f := p.Pages[i]
		path := strings.ReplaceAll(f.CleanPath, "md", "html")
		path = strings.ReplaceAll(path, "./", "")

		_, err = r.WriteString(fmt.Sprintf("* [%s](%s%s)\n", f.Title, p.Path, path))
		if err != nil {
			return fmt.Errorf("write item: %w", err)
		}
	}

	if len(p.Pages) >= 1 {
		_, err = r.WriteString("\n---\n")
		if err != nil {
			return fmt.Errorf("latest title: %w", err)
		}

		latest, err := getFileContent(p.Pages[0])
		if err != nil {
			return fmt.Errorf("getFileContent: %w", err)
		}

		_, err = r.WriteString(latest)
		if err != nil {
			return fmt.Errorf("latest: %w", err)
		}
	} else {
		_, err = r.WriteString("\n[Home](/)")
		if err != nil {
			return fmt.Errorf("home link: %w", err)
		}
	}

	err = r.Close()
	if err != nil {
		return fmt.Errorf("close file: %w", err)
	}

	if fileExists(p.Index) {
		err = os.Remove(p.Index)
		if err != nil {
			return fmt.Errorf("remove old index: %w", err)
		}
	}

	err = os.Rename(p.NewIndex, p.Index)
	if err != nil {
		return fmt.Errorf("move file: %w", err)
	}

	return nil
}
