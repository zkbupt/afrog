package pocs

import (
	"embed"
	"fmt"
	"github.com/zkbupt/afrog/pkg/poc"
	"gopkg.in/yaml.v2"
	"io/fs"
	"strings"
)

//go:embed afrog-pocs/*
var f embed.FS
var EmbedFileList []string

func init() {
	EmbedFileList, _ = EmbedFile()
}

func EmbedFile() ([]string, error) {
	files := []string{}

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// read poc content by name
func EmbedReadContentByName(name string) ([]byte, error) {
	var (
		err    error
		result []byte
	)

	if len(EmbedFileList) == 0 {
		return nil, fmt.Errorf("embed file list is empty")
	}

	for _, file := range EmbedFileList {
		lastSlashIndex := strings.LastIndex(file, "/")
		if lastSlashIndex != -1 {
			fname := file[lastSlashIndex+1:]
			if name == fname || name+".yaml" == fname || name+".yml" == fname {
				fmt.Println(fname)
				return f.ReadFile(file)
			}
		}
	}

	return result, err
}

// read poc struct by path
func EmbedReadPocByPath(path string) (poc.Poc, error) {
	var poc = poc.Poc{}

	file, err := f.Open(path)
	if err != nil {
		return poc, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&poc); err != nil {
		return poc, err
	}
	return poc, nil
}
