package chaotic_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/morhekil/mw/chaotic"
	"github.com/stretchr/testify/require"
)

func TestPolicyApiGet(t *testing.T) {
	c := chaotic.Handler("/chaotic")
	s := httptest.NewServer(c(http.NotFoundHandler()))
	res, err := http.Get(s.URL + "/chaotic/policy")
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	require.Equal(t, `{"Delay":"","DelayP":0}`, string(body))
}

func TestPolicyApiPost(t *testing.T) {
	c := chaotic.Handler("/chaotic")
	s := httptest.NewServer(c(http.NotFoundHandler()))

	res, err := http.Post(s.URL+"/chaotic/policy",
		"application/json",
		strings.NewReader(`{"Delay":"5s","DelayP":0.5}`))
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	require.Equal(t, `{"Delay":"5s","DelayP":0.5}`, string(body))
}

func TestPolicyApiPostWrong(t *testing.T) {
	c := chaotic.Handler("/chaotic")
	s := httptest.NewServer(c(http.NotFoundHandler()))

	res, err := http.Post(s.URL+"/chaotic/policy",
		"application/json",
		strings.NewReader(`{"Delay":"boom","DelayP":0.5}`))
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	require.Equal(t, `{"Delay":"","DelayP":0}`, string(body))
}
