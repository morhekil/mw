package mw

import (
	"net/http"

	"github.com/morhekil/mw/chaotic"
)

/*
Chaotic middleware allows to inject configurable delays and
failures into the requests processed by its underlying HTTP stack.
It provides a configuration page, mounted under the specified URL.

Example:

	http.ListenAndServe(":1234",
		// wrap application handler with chaotic.H,
		// installing its pages under /chaotic URL
		mw.Chaotic("/chaotic")(app),
	)

Read chaotic package documentation for more details: on github
(https://github.com/morhekil/mw/tree/master/chaotic), or on godoc
(http://godoc.org/github.com/morhekil/mw/chaotic).
*/
func Chaotic(url string) func(h http.Handler) http.Handler {
	return chaotic.H(url)
}
