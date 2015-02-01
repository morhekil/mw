package chaotic

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/morhekil/mw/chaotic/bindata"
)

func assets(url string) http.Handler {
	a := &assetfs.AssetFS{
		Asset:    bindata.Asset,
		AssetDir: bindata.AssetDir,
		Prefix:   "",
	}
	return http.StripPrefix(url, http.FileServer(a))
}
