package chaotic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// policyAPI implements HTTP API providing read/write access to policy data
type policyAPI struct {
	p *Policy
}

func (api *policyAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case "POST":
		api.update(w, r)
	default:
		api.show(w, r)
	}

}

func (api *policyAPI) show(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(api.p)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write(j)
	}
}

func (api *policyAPI) update(w http.ResponseWriter, r *http.Request) {
	var np Policy

	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, &np)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		api.p.update(np)
		api.show(w, r)
	}
}
