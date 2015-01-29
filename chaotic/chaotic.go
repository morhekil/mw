package chaotic

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

// Policy describe the desired chaotic behaviour
type Policy struct {
	Delay  string
	DelayP float32
}

// Handler installs its own http routes, and returns
// http.Handler with a potentially chaotic behaviour
func Handler(url string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		p := Policy{}
		mux := http.NewServeMux()

		mux.Handle(url+"/policy", &PolicyAPI{&p})
		mux.Handle(url+"/", assets(url))
		mux.Handle("/", h)

		return mux
	}
}

func assets(url string) http.Handler {
	a := &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "",
	}
	return http.StripPrefix(url, http.FileServer(a))
}
