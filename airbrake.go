package mw

import "net/http"

// AirbrakeNotifier interface describes the common API for error
// notifications
type AirbrakeNotifier interface {
	Notify(interface{}, *http.Request) error
}

// Airbrake middleware
func Airbrake(gb AirbrakeNotifier) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		f := func(hw http.ResponseWriter, hr *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					gb.Notify(r, hr)
					panic(r)
				}
			}()
			h.ServeHTTP(hw, hr)
		}
		return http.HandlerFunc(f)
	}
}
