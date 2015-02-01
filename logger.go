package mw

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger middleware writes all incoming requests into stdout
// using the format:
//     <time>\t<ip>\t"<method> <URI>"\t<runtime>
//
// E.g.
//     2015/02/01 19:19:31 127.0.0.1     "GET /Ping"     31.376µs
//
// IP address prefers X-Real-IP or X-Forwarded-For headers, if present,
// or falls back to http.Request's RemoteAddr otherwise.
func Logger(h http.Handler) http.Handler {
	l := log.New(os.Stdout, "", log.LstdFlags)

	f := func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		h.ServeHTTP(w, r)

		l.Printf("%s\t\"%s %s\"\t%v\n",
			realIP(r),
			r.Method, r.URL.Path,
			time.Since(s),
		)
	}
	return http.HandlerFunc(f)
}
