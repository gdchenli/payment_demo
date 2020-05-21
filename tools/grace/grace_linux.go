package grace

import (
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
)

func Serve(s *http.Server) error {
	return gracehttp.Serve(s)
}
