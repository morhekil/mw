package mw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morhekil/mw"
	"github.com/stretchr/testify/require"
)

func panicHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
}

type testNotifier struct {
	reqs []*http.Request
}

func (tn *testNotifier) Notify(e interface{}, r *http.Request) error {
	tn.reqs = append(tn.reqs, r)
	return nil
}

func TestAirbrakeWithPanic(t *testing.T) {
	nf := testNotifier{
		reqs: make([]*http.Request, 0),
	}

	app := mw.Airbrake(&nf)(panicHandler())
	ht := httptest.NewServer(app)

	hc := http.Client{}
	hc.Get(ht.URL)

	require.Equal(t, 1, len(nf.reqs))
}

func TestAirbrakeWithNoPanic(t *testing.T) {
	nf := testNotifier{
		reqs: make([]*http.Request, 0),
	}

	app := mw.Airbrake(&nf)(http.NotFoundHandler())
	ht := httptest.NewServer(app)

	hc := http.Client{}
	hc.Get(ht.URL)

	require.Equal(t, 0, len(nf.reqs))
}
