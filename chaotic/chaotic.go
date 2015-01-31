package chaotic

import "net/http"

// Handler installs its own http routes, and returns
// http.Handler with a potentially chaotic behaviour
func Handler(url string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		p := policy{}
		p.mux = mux(url, h, p)
		return &p
	}
}

func mux(url string, h http.Handler, p policy) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(url+"/policy", &policyAPI{&p})
	mux.Handle(url+"/", assets(url))
	mux.Handle("/", p.execute(h))

	return mux
}
