package mw_test

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/morhekil/mw"
)

func Example() {
	// Sample application that combines all included
	// middlewares together into a single stack
	app := http.NotFoundHandler()

	hs := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}

	a := alice.New(
		mw.Recover,
		mw.Logger,
		mw.Chaotic("/chaotic"),
		mw.Gzip,
		mw.Headers(hs),
	).Then(app)

	http.ListenAndServe(":1234", a)
}
