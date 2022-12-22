package fs

import (
	"net/http"
	"os"
)

// this fs library is used so that the files in a directory that has been setup as a file server dont get leaked unintentionally

type FileSystem struct {
	Fs http.FileSystem
}

func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
