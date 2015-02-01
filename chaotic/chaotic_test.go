package chaotic_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/morhekil/mw/chaotic"
)

func Example() {
	app := http.NotFoundHandler()
	c := chaotic.H("/chaotic")
	httptest.NewServer(c(app))

	hc := http.Client{}
	hc.Get("/404")     // Hits app's NotFoundHandler
	hc.Get("/chaotic") // Goes to chaotic's configuration page
}
