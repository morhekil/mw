package chaotic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// PolicyAPI implements HTTP API providing read/write access to policy data
type PolicyAPI struct {
	p *Policy
}

func (api *PolicyAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch r.Method {
	case "POST":
		api.update(w, r)
	default:
		api.show(w, r)
	}

}

func (api *PolicyAPI) show(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(api.p)
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write(j)
	}
}

func (api *PolicyAPI) update(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(b, api.p)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		api.show(w, r)
	}
}
