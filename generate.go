package main

import (
	"fmt"
	"os"
	"strings"
)

func writeContent(p *os.File, f []File, title, link string) error {
	if len(f) >= 1 {
		_, err := p.WriteString(fmt.Sprintf("### [%s](/%s)\n", title, link))
		if err != nil {
			return fmt.Errorf("writeContent: %s title: %w", title, err)
		}
		for _, f := range f {
			_, err = p.WriteString(fmt.Sprintf("* [%s](/%s%s)\n", f.Title, link, f.CleanPath))
			if err != nil {
				return fmt.Errorf("writeContent: %s item: %+v, %w", title, f, err)
			}
		}
		_, err = p.WriteString("\n")
		if err != nil {
			return fmt.Errorf("writeContent closer: %w", err)
		}
	}

	return nil
}

func (ic IndexContent) generate() error {
	r, err := os.Create("newIndex.md")
	if err != nil {
		return fmt.Errorf("create newindex: %w", err)
	}

	// Blog
	err = writeContent(r, ic.Blog, "Blog", "blog")
	if err != nil {
		return fmt.Errorf("blog: %w", err)
	}

	// Current
	err = writeContent(r, ic.Current, "Current Projects", "projects/current")
	if err != nil {
		return fmt.Errorf("current projects: %w", err)
	}

	// Past
	err = writeContent(r, ic.Past, "Past Projects", "projects/past")
	if err != nil {
		return fmt.Errorf("past projects: %w", err)
	}

	// Feed
	_, err = r.WriteString("---\n#### Feed\n")
	if err != nil {
		return fmt.Errorf("feed title: %w", err)
	}
	_, err = r.WriteString("[Link](./feed.xml)\n")
	if err != nil {
		return fmt.Errorf("feed link: %w", err)
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

func (p PageContent) generate() error {
	r, err := os.Create(p.NewIndex)
	if err != nil {
		return fmt.Errorf("create index: %v, %w", p.NewIndex, err)
	}

	_, err = r.WriteString("### " + p.Title + "\n")
	if err != nil {
		return fmt.Errorf("write title: %w", err)
	}

	for i := 1; i < len(p.Pages); i++ {
		f := p.Pages[i]
		path := strings.Replace(f.Path, "md", "html", -1)

		_, err = r.WriteString("* [" + f.Title + "](/" + path + ")\n")
		if err != nil {
			return fmt.Errorf("write item: %w", err)
		}
	}

	if len(p.Pages) >= 1 {
		_, err = r.WriteString("\n---\n")
		if err != nil {
			return fmt.Errorf("latest title: %w", err)
		}

		latest, err := getLatest(p.Pages[0])
		if err != nil {
			return fmt.Errorf("getLatest: %w", err)
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

	if _, err = os.Stat(p.Index); !os.IsNotExist(err) {
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
