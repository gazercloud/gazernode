package utilities

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type FileInfo struct {
	Path           string
	Name           string
	Dir            bool
	NameWithoutExt string
	Ext            string
	Size           int64
	Date           time.Time
	Attr           string
}

func (c *FileInfo) SizeAsString() string {
	if c.Dir {
		return "<DIR>"
	}
	div := int64(1)
	uom := ""
	if c.Size >= 1024 {
		div = 1024
		uom = "KB"
	}
	if c.Size >= 1024*1024 {
		div = 1024 * 1024
		uom = "MB"
	}
	if c.Size >= 1024*1024*1024 {
		div = 1024 * 1024 * 1024
		uom = "GB"
	}

	result := strconv.FormatInt(c.Size/div, 10) + " " + uom
	return result
}

func GetDir(path string) ([]FileInfo, error) {
	result := make([]FileInfo, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return result, err
	}

	dirsList := make([]FileInfo, 0)
	filesList := make([]FileInfo, 0)

	for _, f := range files {
		var fileInfo FileInfo
		if f.IsDir() {
			fileInfo.Name = f.Name()
		} else {
			fileInfo.Name = f.Name()
			fileInfo.Ext = filepath.Ext(f.Name())
			if len(fileInfo.Ext) > 0 {
				fileInfo.Ext = fileInfo.Ext[1:]
			}
		}

		fileInfo.NameWithoutExt = strings.TrimSuffix(fileInfo.Name, filepath.Ext(fileInfo.Name))
		fileInfo.Path = path + "/" + f.Name()
		fileInfo.Date = f.ModTime()
		fileInfo.Dir = f.IsDir()
		fileInfo.Size = f.Size()

		if f.IsDir() {
			dirsList = append(dirsList, fileInfo)
		} else {
			filesList = append(filesList, fileInfo)
		}
	}

	sort.Slice(dirsList, func(i, j int) bool {
		return dirsList[i].Name < dirsList[j].Name
	})

	sort.Slice(filesList, func(i, j int) bool {
		return filesList[i].Name < filesList[j].Name
	})

	for _, d := range dirsList {
		result = append(result, d)
	}

	for _, f := range filesList {
		result = append(result, f)
	}

	return result, nil
}
