# chaotic
--
    import "github.com/morhekil/mw/chaotic"

Package chaotic provides stdlib-compatible middleware to inject configurable
delays and failures into the requests processed by its underlying HTTP stack.

It comes with a web page to monitor and configure its behaviour, which also
includes a simple visualisation of requests going through the stack, and
failures introduced.

The minimum viable example of an application with chaotic middleware installed
could be the following:

    app := http.NewServeMux()
    app.Handle("/", http.NotFoundHandler())

    http.ListenAndServe(":1234",
    	// wrap application handler with chaotic.H,
    	// installing its pages under /chaotic URL
    	chaotic.H("/chaotic")(app),
    )

Or it can be cleanly composed with other middleware using alice
(https://github.com/justinas/alice). For example, if we have middlewares called
"logger" and "headers", the full stack can be composed with alice this way:

    	a := alice.New(
    		logger,
    		chaotic.H("/chaotic"),
                    headers,
    	).Then(app)
    	http.ListenAndServe(":1234", a)

Keep in mind, that chaotic will delay or fail middlewares installed after it,
but will not affect middlewares installed earlier - e.g. in this example only
headers middleware is affected by chaotic's behaviour, but logger will run
unaffected every time. This can be used to inject the failure into the required
part of the stack, or even introduce multiple points of failure by mounting
their configuration pages under different URLs.

## Usage

#### func  H

```go
func H(url string) func(h http.Handler) http.Handler
```
H is a net/http handler that installs chaotic's own http routes under the given
base URL, and processes the rest of the stack accordingly to the active policy.

#### type Action

```go
type Action struct {
	Index   int64
	Start   time.Time
	Time    time.Duration
	Delayed bool
	Failed  bool
	Text    string
}
```

Action performed in regards to a single request

#### type Policy

```go
type Policy struct {
	// Public representation of current policy settings
	Delay    string
	DelayP   float32
	FailureP float32
	// Custom function to implement the delay, defaults to time.Sleep.
	DelayFunc func(time.Duration) `json:"-"`
}
```

Policy describes the desired chaotic behaviour

#### func (*Policy) ServeHTTP

```go
func (p *Policy) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

#### func (*Policy) Validate

```go
func (p *Policy) Validate() error
```
Validate should be called to validate public policy data (usually - after a
change), and convert it into internal policy rules, if validation has succeded.
