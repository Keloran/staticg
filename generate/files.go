package generate

import (
  "fmt"
  "io/ioutil"
  "math/rand"
  "os"
  "path/filepath"
  "strings"

  "gopkg.in/yaml.v2"
)

// File ...
type File struct {
  Path      string
  CleanPath string
  Info      os.FileInfo
  Title     string
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

func GetFiles(basePath string) ([]File, error) {
  files := []File{}

  root, err := os.Getwd()
  if err != nil {
    return files, fmt.Errorf("getFiles getwd: %w", err)
  }
  if _, err = os.Stat(root + "/" + basePath); os.IsNotExist(err) {
    return files, fmt.Errorf("%s folder doesn't exist", basePath)
  }

  err = filepath.Walk(root+"/"+basePath, func(path string, info os.FileInfo, err error) error {
    if strings.Contains(info.Name(), ".md") {
      if info.Name() != "index.md" && info.Name() != "newIndex.md" {
        f := File{
          Path:      fmt.Sprintf("%s/%s", basePath, info.Name()),
          CleanPath: fmt.Sprintf("/%s", info.Name()),
          Info:      info,
        }

        title, err := f.getTitle()
        if err != nil {
          return fmt.Errorf("getFiles getTitle: %w", err)
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

func getFileContent(f File) (string, error) {
  dat, err := ioutil.ReadFile(f.Path)
  if err != nil {
    return "", fmt.Errorf("readfile: %w", err)
  }

  return string(dat), nil
}

func (f File) getTitle() (string, error) {
  dat, err := ioutil.ReadFile(f.Path)
  if err != nil {
    return "", fmt.Errorf("getTitle readfile: %w", err)
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
    return "", fmt.Errorf("getTitle cant find end of title")
  }

  return string(dat[startByte:endByte]), nil
}

func fileExists(fn string) bool {
  info, err := os.Stat(fn)
  if os.IsNotExist(err) {
    return false
  }

  return !info.IsDir()
}

func getCV() (string, error) {
  type configStruct struct {
    CV string `yaml:"cv"`
  }

  if !fileExists("_config.yml") {
    return "", fmt.Errorf("no config")
  }

  f := File{
    Path:      "_config.yml",
  }

  cs := configStruct{}
  content, err := getFileContent(f)
  if err != nil {
    return "", fmt.Errorf("getCV content: %w", err)
  }

  err = yaml.Unmarshal([]byte(content), &cs)
  if err != nil {
    return "", fmt.Errorf("getCV unmarshall: %w", err)
  }

  return cs.CV, nil
}

func getOthers() ([]File, error) {
  others := []File{}
  // Feed
  if fileExists("feed.xml") {
    o := File{
      Title: "RSS Feed",
      CleanPath: "/feed.xml",
    }
    others = append(others, o)
  }

  // CV
  cv, err := getCV()
  if err != nil {
    return others, fmt.Errorf("getOthers cv: %w", err)
  }
  if cv != "" {
    o := File{
      CleanPath: cv,
      Title:     "CV",
    }
    others = append(others, o)
  }

  return others, nil
}
