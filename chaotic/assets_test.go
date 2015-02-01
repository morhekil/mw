package chaotic_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morhekil/mw/chaotic"
	"github.com/stretchr/testify/require"
)

func TestAssetsServer(t *testing.T) {
	base := "/chaotic"
	c := chaotic.H(base)
	s := httptest.NewServer(c(http.NotFoundHandler()))

	fs, _ := ioutil.ReadDir("./public")
	for _, f := range fs {
		res, err := http.Get(s.URL + base + "/" + f.Name())
		require.NoError(t, err)
		require.Equal(t, 200, res.StatusCode)
		require.NotEmpty(t, res.ContentLength)
	}
}
