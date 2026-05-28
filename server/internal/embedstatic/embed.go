package embedstatic

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
var assets embed.FS

func ReadFile(path string) ([]byte, error) {
	return assets.ReadFile("assets/" + path)
}

func HasFile(path string) bool {
	_, err := assets.Open("assets/" + path)
	return err == nil
}

func SubFileServer(path string) (http.Handler, error) {
	sub, err := fs.Sub(assets, "assets/"+path)
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(sub)), nil
}
