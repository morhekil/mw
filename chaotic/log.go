package chaotic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const hist = 250

// Action performed in regards to a single request
type Action struct {
	Index   int64
	Time    time.Duration
	Delayed bool
	Failed  bool
	Text    string
}

// Log of all requests processed
type Log struct {
	count int64
	items [hist]Action
}

// Push another action into the log
func (l *Log) Push(a Action) {
	a.Index = l.count
	l.items[l.count%hist] = a
	l.count++
}

// Export current log data as JSON array
func (l *Log) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	min := l.count - hist
	if min < 0 {
		min = 0
	}
	res := make([]string, l.count-min)

	for i := min; i < l.count; i++ {
		s, err := json.Marshal(l.items[(i+hist)%hist])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		res[i-min] = string(s)
	}
	fmt.Fprintf(w, "[%s]", strings.Join(res, ",\n"))
}
