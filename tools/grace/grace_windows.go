package grace

import (
	"net/http"
)

func Serve(s *http.Server) error {
	// The code will be a bit more complex when you add support
	// of Serve(s ...*http.Server).
	return s.ListenAndServe()
}
