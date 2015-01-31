package chaotic_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/morhekil/mw/chaotic"
	"github.com/stretchr/testify/require"
)

func TestDelayValuePolicy(t *testing.T) {
	c := 0
	d5s, _ := time.ParseDuration("5s")
	f := func(d time.Duration) {
		require.Equal(t, d5s, d)
		c++
	}
	p := chaotic.Policy{
		Delay:     "5s",
		DelayP:    1.0,
		DelayFunc: f,
	}
	require.NoError(t, p.Validate())

	s := httptest.NewServer(&p)
	for i := 0; i < 10; i++ {
		res, err := http.Get(s.URL + "/ping")
		require.NoError(t, err)
		require.Equal(t, 404, res.StatusCode)
	}
	require.Equal(t, 10, c)
}

func TestDelayChancePolicy(t *testing.T) {
	c := 0
	d5s, _ := time.ParseDuration("5s")
	f := func(d time.Duration) {
		require.Equal(t, d5s, d)
		c++
	}
	p := chaotic.Policy{
		Delay:     "5s",
		DelayP:    0.5,
		DelayFunc: f,
	}
	require.NoError(t, p.Validate())

	s := httptest.NewServer(&p)
	for i := 0; i < 30; i++ {
		res, err := http.Get(s.URL + "/ping")
		require.NoError(t, err)
		require.Equal(t, 404, res.StatusCode)
	}
	// with 50% delay probability, our delay function should've been
	// called 15 +/- 13 times
	require.InDelta(t, 15, c, 13)
}

func TestDelayNativeFunc(t *testing.T) {
	p := chaotic.Policy{
		Delay:  "300ms",
		DelayP: 0.5,
	}
	require.NoError(t, p.Validate())

	s := httptest.NewServer(&p)
	defer func(st time.Time) {
		require.True(t, time.Since(st).Nanoseconds() > 300e6,
			"execution should take at least 300ms")
	}(time.Now())

	res, err := http.Get(s.URL + "/ping")
	require.NoError(t, err)
	require.Equal(t, 404, res.StatusCode)
}
