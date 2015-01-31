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

func resBody(t *testing.T, res *http.Response) string {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	return string(body)
}

func testServer() *httptest.Server {
	c := chaotic.Handler("/chaotic")
	s := httptest.NewServer(c(http.NotFoundHandler()))
	return s
}

func TestPolicyApiGet(t *testing.T) {
	res, err := http.Get(testServer().URL + "/chaotic/policy")
	require.NoError(t, err)
	require.Equal(t, `{"Delay":"","DelayP":0,"FailureP":0}`, resBody(t, res))
}

func TestPolicyApiPost(t *testing.T) {
	res, err := http.Post(testServer().URL+"/chaotic/policy",
		"application/json",
		strings.NewReader(`{"Delay":"5s","DelayP":0.5,"FailureP":0.3}`))
	require.NoError(t, err)

	require.Equal(t, `{"Delay":"5s","DelayP":0.5,"FailureP":0.3}`,
		resBody(t, res))
}

func TestPolicyApiPostWrong(t *testing.T) {
	res, err := http.Post(testServer().URL+"/chaotic/policy",
		"application/json",
		strings.NewReader(`{"Delay":"boom","DelayP":0.5,"FailureP":0.3}`))
	require.NoError(t, err)

	require.Equal(t, `{"Delay":"","DelayP":0,"FailureP":0}`,
		resBody(t, res))
}

func TestPolicyApiPostMalformed(t *testing.T) {
	res, err := http.Post(testServer().URL+"/chaotic/policy",
		"application/json",
		strings.NewReader(`error`))
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, 500)
}
