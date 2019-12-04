package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EditorLink struct {
	Title []byte
	Link  []byte
}

type EditorContent struct {
	Blog      bool
	Title     []byte
	Content   []byte
	PageLinks []EditorLink
}

func getResponse(delim byte) ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadBytes(delim)
	if err != nil {
		return []byte{}, fmt.Errorf("readstring: %w", err)
	}

	return input[0 : len(input)-1], nil
}

func editor(root string) error {
	ec, err := termEditor()
	if err != nil {
		return fmt.Errorf("termeditor: %w", err)
	}

	filename := strings.ToLower(strings.Replace(string(ec.Title), " ", "_", -1))
	path := root + "projects/current"
	if ec.Blog {
		path = root + "blog"
	}

	f, err := os.Create(fmt.Sprintf("%s/%s.md", path, filename))
	if err != nil {
		return fmt.Errorf("file create: %w", err)
	}

	_, err = f.WriteString("### " + string(ec.Title) + "\n")
	if err != nil {
		return fmt.Errorf("title: %w", err)
	}

	_, err = f.Write(ec.Content)
	if err != nil {
		return fmt.Errorf("content: %w", err)
	}

	_, err = f.WriteString("---\n\n")
	if err != nil {
		return fmt.Errorf("link start: %w", err)
	}
	for _, link := range ec.PageLinks {
		_, err := f.WriteString("[" + string(link.Title) + "](" + string(link.Link) + ")\n--\n")
		if err != nil {
			return fmt.Errorf("link write: %w", err)
		}
	}

	// Home link
	_, err = f.WriteString("[Home](/)\n")
	if err != nil {
		return fmt.Errorf("home: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func termEditor() (EditorContent, error) {
	ec := EditorContent{}

	fmt.Printf("Project or Blog q/b ? ")
	t, err := getResponse('\n')
	if err != nil {
		return ec, fmt.Errorf("question err: %w", err)
	}
	ts := string(t)
	switch ts {
	case "blog":
	case "b":
		ec.Blog = true
	default:
		ec.Blog = false
	}

	fmt.Printf("Title: ")
	title, err := getResponse('\n')
	if err != nil {
		return ec, fmt.Errorf("title: %w", err)
	}
	ec.Title = title

	fmt.Printf("Content (tilde to quit): ")
	content, err := getResponse('~')
	if err != nil {
		return ec, fmt.Errorf("content: %w", err)
	}
	ec.Content = content

	fmt.Printf("How many links: ")
	linkCount, err := getResponse('\n')
	if err != nil {
		return ec, fmt.Errorf("link count: %w", err)
	}
	linkCountI, err := strconv.Atoi(string(linkCount))
	if err != nil {
		return ec, fmt.Errorf("convert link count: %w", err)
	}
	els := []EditorLink{}
	for i := 0; i < linkCountI; i++ {
		el := EditorLink{}

		fmt.Printf("Link Title: ")
		t, err := getResponse('\n')
		if err != nil {
			return ec, fmt.Errorf("editor link title: %w", err)
		}
		el.Title = t

		fmt.Printf("Link: ")
		l, err := getResponse('\n')
		if err != nil {
			return ec, fmt.Errorf("editor link link: %w", err)
		}
		el.Link = l

		els = append(els, el)
	}
	ec.PageLinks = els

	fmt.Printf("\nSave y/n/q ? ")
	save, err := getResponse('\n')
	if err != nil {
		return ec, fmt.Errorf("save: %w", err)
	}

	switch save[0] {
	case YES:
		return ec, nil
	case QUIT:
		return EditorContent{}, nil
	case NO:
		return termEditor()
	}

	return ec, nil
}
