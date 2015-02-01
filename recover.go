package mw

import (
	"log"
	"net/http"
	"os"
)

// Recover from panic when it happens in the application.
// Details of the panic are written into Stderr, together with
// some information about the request- client's real IP address,
// HTTP method and URL path.
func Recover(h http.Handler) http.Handler {
	l := log.New(os.Stderr, "", log.LstdFlags)

	errf := func(r *http.Request) {
		if e := recover(); e != nil {
			l.Printf("%s\t\"%s %s\"\tPANIC: %v\n",
				realIP(r), r.Method, r.URL.Path, e,
			)
		}
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		defer errf(r)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
