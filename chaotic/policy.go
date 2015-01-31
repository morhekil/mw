package chaotic

import (
	"math/rand"
	"net/http"
	"time"
)

// Policy describes the desired chaotic behaviour
type Policy struct {
	// Public representation of current policy settings
	Delay    string
	DelayP   float32
	FailureP float32
	// Custom function to implement the delay, defaults to time.Sleep.
	DelayFunc func(time.Duration) `json:"-"`
	// converted value of Delay
	delay time.Duration
	// ServeMux to serve this policy as http middleware
	mux http.Handler
}

func (p *Policy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var next http.HandlerFunc
	if p.mux != nil {
		next = p.mux.ServeHTTP
	} else {
		next = http.NotFound
	}
	h := p.execute(next)
	h(w, r)
}

// Update policy with a new set of parameters.
// If update fails (e.g. Delay value is misformatted, etc),
// the existing policy will remain intact, and an error is returned
func (p *Policy) update(np Policy) error {
	if err := np.Validate(); err != nil {
		return err
	}

	p.Delay = np.Delay
	p.DelayP = np.DelayP
	p.FailureP = np.FailureP
	p.delay = np.delay

	return nil
}

func (p *Policy) execute(hf http.HandlerFunc) http.HandlerFunc {
	if p.delay != 0 && rand.Float32() < p.DelayP {
		p.execDelay()
	}
	if rand.Float32() < p.FailureP {
		return p.execFailure()
	}

	return hf
}

// Validate should be called to validate public policy data
// (usually - after a change), and convert it into internal policy
// rules, if validation has succeded.
func (p *Policy) Validate() error {
	var (
		d   time.Duration
		err error
	)

	if p.Delay != "" {
		d, err = time.ParseDuration(p.Delay)
	}
	if err != nil {
		return err
	}

	p.delay = d
	return nil
}

// Execute delay according to the current policy (delay function and value)
func (p *Policy) execDelay() {
	if p.DelayFunc != nil {
		p.DelayFunc(p.delay)
	} else {
		time.Sleep(p.delay)
	}
}

func (p *Policy) execFailure() http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I'm an agent of chaos", 500)
	}
	return http.HandlerFunc(h)
}
