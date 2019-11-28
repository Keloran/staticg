package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path string
	Info os.FileInfo
	Title string
}

type IndexContent struct {
	Blog []File
	Current []File
	Past []File
}

type PageContent struct {
	Title string
	NewIndex string
	Index string
	Pages []File
}

func main() {
	ic := IndexContent{}

	pages, err := blogPages()
	if err != nil {
		fmt.Printf("blog err: %v\n", err)
		return
	}
	ic.Blog = pages

	pages, err = currentProjects()
	if err != nil {
		fmt.Printf("current projects: %v\n", err)
		return
	}
	ic.Current = pages

	pages, err = pastProjects()
	if err != nil {
		fmt.Printf("past projects: %v\n", err)
		return
	}
	ic.Past = pages

	err = ic.generate()
	if err != nil {
		fmt.Printf("generate index: %v\n", err)
		return
	}

	fmt.Printf("all files created\n")
	return
}

func blogPages() ([]File, error) {
	files, err := getFiles("blog")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	} 

	p := PageContent{
		Title: "Blog",
		NewIndex: "blog/newIndex.md",
		Index: "blog/index.md",
		Pages: files,
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

func currentProjects() ([]File, error) {
	files, err := getFiles("projects/current")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title: "Current Projects",
		NewIndex: "projects/current/newIndex.md",
		Index: "projects/current/index.md",
		Pages: files,
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

func pastProjects() ([]File, error) {
	files, err := getFiles("projects/past")
	if err != nil {
		return []File{}, fmt.Errorf("getFiles: %w", err)
	}

	p := PageContent{
		Title: "Past Projects",
		NewIndex: "projects/past/newIndex.md",
		Index: "projects/past/index.md",
		Pages: files,
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

func getFiles(basePath string) ([]File, error) {
	files := []File{}
	
	root, err := os.Getwd()
	if err != nil {
		return files, fmt.Errorf("getwd: %w", err)
	}
	if _, err = os.Stat(root + "/" + basePath); os.IsNotExist(err) {
		return files, fmt.Errorf("%s folder doesn't exist", basePath)
	}

	err = filepath.Walk(root + "/" + basePath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".md") {
			if info.Name() != "index.md" {
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

func (f File)getTitle() (string, error) {
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
			startByte = i+2
		}
		if b == newLine {
			if endByte == 0 {
				endByte = i
			}
		}
	}

	if endByte == 0 {
		return "", fmt.Errorf("cant find end of title")
	}

	return string(dat[startByte:endByte]), nil
}

func getLatest(f File) (string, error) {
	dat, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return "", fmt.Errorf("readfile: %w", err)
	}

	return string(dat), nil
}

func (ic IndexContent)generate() error {
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

func (p PageContent)generate() error {
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
		_, err = r.WriteString("* [" + f.Title + "](/" + f.Path + ")\n")
		if err != nil {
			return fmt.Errorf("write item: %w", err)
		}
	}

	if len(p.Pages) >= 1 {
		_, err = r.WriteString("---\n")
		if err != nil {
			return fmt.Errorf("latest title: %w", err)
		}

		latest, err := getLatest(p.Pages[0])
		if err != nil {
			return fmt.Errorf("getLatest: %w", err)
		}

		_, err = r.WriteString(latest)
		if err != nil {
			return fmt.Errorf("write latest: %w", err)
		}
	}

//	_, err = r.WriteString("---\n[Home](/)\n")
//	if err != nil {
//		return fmt.Errorf("write home link: %w", err)
//	}

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

func sortFiles(files []File) []File {
	ret := []File{}

	for _, f := range files {
		for _, ff := range files {
			if ff.Info.ModTime().Unix() >= f.Info.ModTime().Unix() {
				o := ret
				t := []File{}
				t = append(t, ff)
				for _, x := range o {
					if x.Path != ff.Path {
						t = append(t, x)
					}
				}
				ret = t
			}
		}
	}

	return ret
}
