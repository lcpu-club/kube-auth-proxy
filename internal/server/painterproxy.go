package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (s *Server) initPainterProxy() {
	urlStr := "http://painter.painter-system"
	urlURL, _ := url.Parse(urlStr)

	painterRev := httputil.NewSingleHostReverseProxy(urlURL)

	s.mux.HandleFunc("/_/painter/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/_/painter")
		painterRev.ServeHTTP(w, r)
	})
}
