package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path  string
	Info  os.FileInfo
	Title string
}

type IndexContent struct {
	Blog    []File
	Current []File
	Past    []File
}

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

const YES = 121
const NO = 110
const QUIT = 113

func _main(args []string) error {
	root := "./"
	if len(args) >= 1 {
		if args[0] != "" {
			root = args[0] + "/"
		}
	}

	ic := IndexContent{}

	fmt.Printf("Create new item y/n ? ")
	create, err := getResponse('\n')
	if create[0] == YES {
		err := editor(root)
		if err != nil {
			return fmt.Errorf("editor err: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("item question: %w", err)
	}

	pages, err := blogPages(root)
	if err != nil {
		return fmt.Errorf("blog err: %w", err)
	}
	ic.Blog = pages

	pages, err = currentProjects(root)
	if err != nil {
		return fmt.Errorf("current projects: %w", err)
	}
	ic.Current = pages

	pages, err = pastProjects(root)
	if err != nil {
		return fmt.Errorf("past projects: %w", err)
	}
	ic.Past = pages

	err = ic.generate()
	if err != nil {
		return fmt.Errorf("generate index: %w", err)
	}

	err = ic.generateFeed()
	if err != nil {
		return fmt.Errorf("generate feed: %w", err)
	}

	fmt.Printf("all files created\n")
	return nil
}

func sortFiles(files []File) []File {
	if len(files) < 2 {
		return files
	}

	left, right := 0, len(files)-1
	pivot := rand.Int() % len(files)
	files[pivot], files[right] = files[right], files[pivot]

	for i, _ := range files {
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
					Path: basePath + "/" + info.Name(),
					Info: info,
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
