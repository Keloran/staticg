package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// File ...
type File struct {
	Path      string
	CleanPath string
	Info      os.FileInfo
	Title     string
}

// IndexContent ...
type IndexContent struct {
	Blog    []File
	Current []File
	Past    []File
}

// PageContent ...
type PageContent struct {
	Title    string
	NewIndex string
	Index    string
	Pages    []File
}

func main() {
	err := _main(os.Args[1:])
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		return
	}
}

// YES keycode for y
const YES = 121

// NO keycode for n
const NO = 110

// QUIT keycode for q
const QUIT = 113

func _main(args []string) error {
	root := "./"
	if len(args) >= 1 {
		if args[0] != "" {
			root = args[0] + "/"
		}
	}

	ic := IndexContent{}

	errChan := make(chan error)
	generated := make(chan bool)

	fmt.Printf("Create new item y/n ? ")
	create, err := getResponse('\n')
	if create[0] == YES {
		err := editor(root)
		if err != nil {
			errChan <- fmt.Errorf("editor err: %w", err)
		}
	}
	if err != nil {
		errChan <- fmt.Errorf("item question: %w", err)
	}

	pages, err := blogPages(root)
	if err != nil {
		errChan <- fmt.Errorf("blog err: %w", err)
	}
	ic.Blog = pages
	fmt.Printf("blogs gathered\n")

	pages, err = currentProjects(root)
	if err != nil {
		errChan <- fmt.Errorf("current projects: %w", err)
	}
	ic.Current = pages
	fmt.Printf("current projects gathered\n")

	pages, err = pastProjects(root)
	if err != nil {
		errChan <- fmt.Errorf("past projects: %w", err)
	}
	ic.Past = pages
	fmt.Printf("past projects gathered\n")

	go func() {
		err = ic.generate()
		if err != nil {
			errChan <- fmt.Errorf("generate index: %w", err)
		}
		fmt.Printf("index pages generated\n")

		err = ic.generateFeed()
		if err != nil {
			errChan <- fmt.Errorf("generate feed: %w", err)
		}
		fmt.Printf("rss generated\n")
		generated <- true
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("main err: %w", err)
	case <-generated:
		return nil
	}
}

func sortFiles(files []File) []File {
	if len(files) < 2 {
		return files
	}

	left, right := 0, len(files)-1
	pivot := rand.Int() % len(files)
	files[pivot], files[right] = files[right], files[pivot]

	for i := range files {
		if files[i].Info.ModTime().Unix() > files[right].Info.ModTime().Unix() {
			files[left], files[i] = files[i], files[left]
			left++
		}
	}
	files[left], files[right] = files[right], files[left]
	sortFiles(files[:left])
	sortFiles(files[left+1:])

	return files
}

func getFiles(basePath string) ([]File, error) {
	files := []File{}

	root, err := os.Getwd()
	if err != nil {
		return files, fmt.Errorf("getwd: %w", err)
	}
	if _, err = os.Stat(root + "/" + basePath); os.IsNotExist(err) {
		return files, fmt.Errorf("%s folder doesn't exist", basePath)
	}

	err = filepath.Walk(root+"/"+basePath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".md") {
			if info.Name() != "index.md" && info.Name() != "newIndex.md" {
				f := File{
					Path:      basePath + "/" + info.Name(),
					CleanPath: "/" + info.Name(),
					Info:      info,
				}

				title, err := f.getTitle()
				if err != nil {
					return fmt.Errorf("getTitle: %w", err)
				}
				f.Title = title

				files = append(files, f)
			}
		}
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("walk: %w", err)
	}

	return sortFiles(files), nil
}

func getLatest(f File) (string, error) {
	dat, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return "", fmt.Errorf("readfile: %w", err)
	}

	return string(dat), nil
}

func (f File) getTitle() (string, error) {
	dat, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return "", fmt.Errorf("readfile: %w", err)
	}

	hashBytes := []byte{35, 32}
	newLine := byte(10)

	startByte := 0
	endByte := 0

	for i, b := range dat {
		if b == hashBytes[0] && dat[i+1] == hashBytes[1] {
			startByte = i + 2
		}
		if b == newLine {
			if endByte == 0 {
				endByte = i
			}
		}

		if endByte != 0 && startByte != 0 {
			break
		}
	}

	if endByte == 0 {
		return "", fmt.Errorf("cant find end of title")
	}

	return string(dat[startByte:endByte]), nil
}
