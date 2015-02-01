package chaotic

import (
	"fmt"
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
	// Log of processed actions
	logger *logger
	// converted value of Delay
	delay time.Duration
	// next to serve this policy as http middleware
	next http.Handler
	// channel for log messages
	ch chan Action
}

func (p *Policy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hf := p.upstream()
	a := Action{
		Start: time.Now(),
	}

	if p.delay != 0 && rand.Float32() < p.DelayP {
		a.Delayed = true
		p.execDelay()
	}
	if rand.Float32() < p.FailureP {
		a.Failed = true
		hf = p.execFailure()
	}

	a.Time = time.Since(a.Start)
	a.Text = fmt.Sprintf("%s %s", r.Method, r.URL)
	if p.ch != nil {
		p.ch <- a
	}
	hf(w, r)
}

func (p *Policy) upstream() http.HandlerFunc {
	if p.next != nil {
		return p.next.ServeHTTP
	}
	return http.NotFound
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
