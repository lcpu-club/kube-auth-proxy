package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (s *Server) initPainterProxy() {
	urlStr := "http://painter.painter-system:23333"
	urlURL, _ := url.Parse(urlStr)

	painterRev := httputil.NewSingleHostReverseProxy(urlURL)

	painterHandler := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/_/painter")
		painterRev.ServeHTTP(w, r)
	}
	s.mux.HandleFunc("/_/painter", painterHandler)
	s.mux.HandleFunc("/_/painter/", painterHandler)
}
