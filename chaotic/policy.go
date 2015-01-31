package chaotic

import (
	"net/http"
	"time"
)

// policy describes the desired chaotic behaviour
type policy struct {
	Delay  string
	DelayP float32
	delay  time.Duration
	mux    http.Handler
}

func (p *policy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mux.ServeHTTP(w, r)
}

// Update policy with a new set of parameters.
// If update fails (e.g. Delay value is misformatted, etc),
// the existing policy will remain intact, and an error is returned
func (p *policy) update(np policy) error {
	if err := np.validate(); err != nil {
		return err
	}

	p.Delay = np.Delay
	p.DelayP = np.DelayP
	p.delay = np.delay

	return nil
}

func (p *policy) execute(h http.Handler) http.Handler {
	return h
}

func (p *policy) validate() error {
	d, err := time.ParseDuration(p.Delay)
	if err != nil {
		return err
	}

	p.delay = d
	return nil
}
