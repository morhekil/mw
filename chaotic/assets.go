package chaotic

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

func assets(url string) http.Handler {
	a := &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "",
	}
	return http.StripPrefix(url, http.FileServer(a))
}
