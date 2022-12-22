package frontend

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"strings"

	psFs "PasswordServer2/lib/fs"
)

//go:generate yarn
//go:generate yarn run build
//go:embed all:build
var files embed.FS

func SvelteKitHandler(path string) http.Handler {
	fsys, err := fs.Sub(files, "build")

	if err != nil {
		panic(err)
	}
	filesystem := psFs.FileSystem{Fs: http.FS(fsys)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, path)

		_, err := filesystem.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			path += ".html"
		}

		r.URL.Path = path
		http.FileServer(filesystem).ServeHTTP(w, r)
	})
}
