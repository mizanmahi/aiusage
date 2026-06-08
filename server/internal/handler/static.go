package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StaticHandler struct {
	dir        string
	fileServer http.Handler
}

func NewStaticHandler(dir string) (*StaticHandler, error) {
	indexPath := filepath.Join(dir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		return nil, err
	}

	return &StaticHandler{
		dir:        dir,
		fileServer: http.FileServer(http.Dir(dir)),
	}, nil
}

func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.NotFound(w, r)
		return
	}

	path := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if path == "." {
		h.serveIndex(w, r)
		return
	}
	if !filepath.IsLocal(path) {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(h.dir, path)
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		h.fileServer.ServeHTTP(w, r)
		return
	}

	h.serveIndex(w, r)
}

func (h *StaticHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(h.dir, "index.html"))
}
