# chaotic

```go
import "github.com/morhekil/mw/chaotic"
```

Package `chaotic` provides stdlib-compatible middleware to inject configurable
delays and failures into the requests processed by its underlying HTTP stack.

It comes with a web page to monitor and configure its behaviour, which also
includes a simple visualisation of requests going through the stack, and
failures introduced.

The minimum viable example of an application with chaotic middleware installed
could be the following:

```go
app := http.NewServeMux()
app.Handle("/", http.NotFoundHandler())

http.ListenAndServe(":1234",
    // wrap application handler with chaotic.H,
    // installing its pages under /chaotic URL
    chaotic.H("/chaotic")(app),
)
```

Or it can be cleanly composed with other middleware using
[alice](https://github.com/justinas/alice).  For example, if we have
middlewares called `logger` and `headers`, the full stack can be
composed with alice this way:

```go
a := alice.New(
	logger,
	chaotic.H("/chaotic"),
	    headers,
).Then(app)
http.ListenAndServe(":1234", a)
```

Keep in mind, that chaotic will delay or fail middlewares installed after it,
but will not affect middlewares installed earlier - e.g. in this example only
`headers` middleware is affected by chaotic's behaviour, but `logger` will run
unaffected every time. This can be used to inject the failure into the required
part of the stack, or even introduce multiple points of failure by mounting
their configuration pages under different URLs.

Full API documentation is available at [godoc](http://godoc.org/github.com/morhekil/mw/chaotic).

## Configuration page

When installed in the middleware stack, chaotic servers its own configuration
page under its URL. Of course, the use of this page is completely optional,
and you can achieve the same result by just talking to chaotic's API directly.

Below is an example of the configuration page in action:

![demo](http://f.falsum.me/image/411m1c0a2r04/chaotic.gif)