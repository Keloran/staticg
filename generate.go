package main

import (
	"fmt"
	"os"
	"strings"
)

func (ic IndexContent) generate() error {
	r, err := os.Create("newIndex.md")
	if err != nil {
		return fmt.Errorf("create newindex: %w", err)
	}

	// Blog
	_, err = r.WriteString("### [Blog](/blog)\n")
	if err != nil {
		return fmt.Errorf("blog title: %w", err)
	}
	for _, f := range ic.Blog {
		_, err = r.WriteString("* [" + f.Title + "](/" + f.Path + ")\n")
		if err != nil {
			return fmt.Errorf("blog item: %+v, %w", f, err)
		}
	}

	// Current
	_, err = r.WriteString("\n\n### [Current Projects](/projects/current)\n")
	if err != nil {
		return fmt.Errorf("current title: %w", err)
	}
	for _, f := range ic.Current {
		_, err = r.WriteString("* [" + f.Title + "](/" + f.Path + ")\n")
		if err != nil {
			return fmt.Errorf("current item: %+v, %w", f, err)
		}
	}

	// Past
	_, err = r.WriteString("\n\n### [Past Projects](/projects/past)\n")
	if err != nil {
		return fmt.Errorf("past title: %w", err)
	}
	for _, f := range ic.Past {
		_, err = r.WriteString("* [" + f.Title + "](/" + f.Path + ")\n")
		if err != nil {
			return fmt.Errorf("past time: %+v, %w", f, err)
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

	err = os.Rename(p.NewIndex, p.Index)
	if err != nil {
		return fmt.Errorf("move file: %w", err)
	}

	return nil
}
