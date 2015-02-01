package chaotic_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/morhekil/mw/chaotic"
	"github.com/stretchr/testify/require"
)

func TestLogExport(t *testing.T) {
	s := testServer()
	for i := 0; i < 5; i++ {
		_, err := http.Get(s.URL + "/ping")
		require.NoError(t, err)
	}

	res, err := http.Get(s.URL + "/chaotic/log")
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)

	l := make([]chaotic.Action, 5)
	b := resBody(t, res)
	err = json.Unmarshal([]byte(b), &l)
	require.NoError(t, err)
	require.Equal(t, 5, len(l))

	for i := 0; i < 5; i++ {
		require.NotEqual(t, chaotic.Action{}, l[i])
	}
}

func TestLogClear(t *testing.T) {
	s := testServer()
	for i := 0; i < 5; i++ {
		_, err := http.Get(s.URL + "/ping")
		require.NoError(t, err)
	}

	res, err := http.Post(s.URL+"/chaotic/log", "application/json",
		strings.NewReader(""))
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)

	res, err = http.Get(s.URL + "/chaotic/log")
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)

	require.Equal(t, "[]", resBody(t, res))
}
