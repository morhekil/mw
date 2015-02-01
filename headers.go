package mw

import "net/http"

// Headers to be written into all responses from the application.
// Example:
//     hs := map[string]string{
//         "Content-Type":   "application/json; charset=utf-8",
//     }
//     http.ListenAndServe(":1234", mw.Headers(hs)(app))
//
func Headers(hs map[string]string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			for k, v := range hs {
				w.Header().Set(k, v)
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}
