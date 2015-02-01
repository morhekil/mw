package chaotic

import "net/http"

// Handler installs its own http routes, and returns
// http.Handler with a potentially chaotic behaviour
func Handler(url string) func(h http.Handler) http.Handler {
	p := Policy{
		ch:  make(chan Action),
		Log: &Log{},
	}
	go p.Log.Pull(p.ch)
	return mux(url, p)
}

func mux(url string, p Policy) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		p.next = h

		mux := http.NewServeMux()
		mux.Handle(url+"/policy", &policyAPI{&p})
		mux.Handle(url+"/log", p.Log)
		mux.Handle(url+"/", assets(url))
		mux.Handle("/", &p)

		return mux
	}
}
