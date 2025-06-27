package identities

import "net/http"

type httpServer struct{}

func newServeMux(s *httpServer) *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func (h *httpServer) wellKnown(w http.ResponseWriter, r *http.Request) {}

func (h *httpServer) authorization(w http.ResponseWriter, r *http.Request) {}
