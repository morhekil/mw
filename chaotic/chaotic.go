/*
Package chaotic provides stdlib-compatible middleware to inject
configurable delays and failures into the requests processed by its
underlying HTTP stack.

It comes with a web page to monitor and configure its behaviour, which
also includes a simple visualisation of requests going through the
stack, and failures introduced.

An example of a configuration page in action can be found on GitHub
(https://github.com/morhekil/mw/tree/master/chaotic#chaotic).

The minimum viable example of an application with chaotic middleware
installed could be the following:

	app := http.NewServeMux()
	app.Handle("/", http.NotFoundHandler())

	http.ListenAndServe(":1234",
		// wrap application handler with chaotic.H,
		// installing its pages under /chaotic URL
		mw.Chaotic("/chaotic")(app),
	)

Or it can be cleanly composed with other middleware using alice
(https://github.com/justinas/alice). For example,
if we have middlewares called "mw.Logger" and "mw.Recover", the full stack
can be composed with alice this way:

	a := alice.New(
		mw.Logger,
		mw.Chaotic("/chaotic"),
                mw.Recover,
	).Then(app)
	http.ListenAndServe(":1234", a)

Keep in mind, that chaotic will delay or fail middlewares installed after it,
but will not affect middlewares installed earlier - e.g. in this example only
mw.Recover middleware is affected by chaotic's behaviour, but mw.Logger will run
unaffected every time. This can be used to inject the failure into the required
part of the stack, or even introduce multiple points of failure by mounting
their configuration pages under different URLs.
*/
package chaotic

import "net/http"

// H is a net/http handler that installs chaotic's own http routes
// under the given base URL, and processes the rest of the stack
// accordingly to the active policy.
func H(url string) func(h http.Handler) http.Handler {
	p := Policy{
		ch:     make(chan Action),
		logger: &logger{},
	}
	go p.logger.Pull(p.ch)
	return mux(url, p)
}

func mux(url string, p Policy) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		p.next = h

		mux := http.NewServeMux()
		mux.Handle(url+"/policy", &policyAPI{&p})
		mux.Handle(url+"/log", p.logger)
		mux.Handle(url+"/", assets(url))
		mux.Handle("/", &p)

		return mux
	}
}
