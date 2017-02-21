package zip

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Compress interface {
	Close() error
	WriteHead(path string, info os.FileInfo) error
	Write(p []byte) (int, error)
}

func walk(path string, compresser Compress) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(path)
	}

	filepath.Walk(path, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		root = filepath.ToSlash(root)
		fileroot := root
		if root = strings.TrimPrefix(root, path); root == "" {
			root = baseDir
		} else {
			root = baseDir + root
		}
		err = compresser.WriteHead(root, info)
		if err != nil {
			return nil
		}
		F, err := os.Open(fileroot)
		if err != nil {
			return nil
		}
		io.Copy(compresser, F)
		F.Close()
		return nil
	})
	return nil
}
