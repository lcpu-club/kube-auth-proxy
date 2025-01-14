package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/lcpu-club/kube-auth-proxy/internal/utils"
)

func (s *Server) generateToken(uid string) string {
	return utils.GenerateRandomString(*s.conf.TokenLength, "sk:"+uid+":")
}

func (s *Server) getToken(token string) (*ImpersonateInfo, error) {
	iiStr, err := s.stor.Load(token)
	if err != nil {
		return nil, err
	}

	return ImpersonateInfoFromString(iiStr)
}

func (s *Server) setToken(token string, ii *ImpersonateInfo) error {
	return s.stor.Store(token, ii.String(), *s.conf.TokenExpiration)
}

func (s *Server) deleteToken(token string) error {
	return s.stor.Delete(token)
}

func (s *Server) listTokens(uid string) ([]string, error) {
	prefix := "sk:"
	if uid != "" {
		prefix += uid + ":"
	}
	return s.stor.List(prefix)
}

func (s *Server) initToken() {
	s.kubeconfigTemplate = template.Must(
		template.ParseFiles(*s.conf.KubeconfigTemplatePath),
	)

	s.mux.HandleFunc("/_/tokens", s.handleToken)
	s.mux.HandleFunc("/_/tokens/", s.handleToken)
}

func (s *Server) handleToken(w http.ResponseWriter, r *http.Request) {
	ii, err := s.authenticate(r)
	if err != nil {
		http.Error(w, "Failed to authenticate", http.StatusUnauthorized)
	}
	uid := ii.UID
	if uid == "" {
		uid = ii.Username
	}
	if uid == "" {
		http.Error(w, "Failed to get UID", http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		s.handleGetToken(w, r, uid)
	case http.MethodPost:
		s.handlePostToken(w, ii, uid)
	case http.MethodDelete:
		s.handleDeleteToken(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetTokenKubeconfig(w http.ResponseWriter, r *http.Request, uid string) {
	token := r.URL.Path[len("/_/tokens"):]
	if token == "" || token == "/" {
		http.Error(w, "No token provided", http.StatusBadRequest)
	}

	token = token[1:] // Remove leading slash
	token = strings.TrimSuffix(token, "/kubeconfig")

	vars := struct {
		Username string
		Token    string
	}{
		Username: uid,
		Token:    token,
	}

	err := s.kubeconfigTemplate.Execute(w, vars)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to render kubeconfig", http.StatusInternalServerError)
	}
}

func (s *Server) handleGetToken(w http.ResponseWriter, r *http.Request, uid string) {
	if strings.HasSuffix(r.URL.Path, "/kubeconfig") {
		s.handleGetTokenKubeconfig(w, r, uid)
		return
	}

	tokens, err := s.listTokens(uid)
	if err != nil {
		http.Error(w, "Failed to list tokens", http.StatusInternalServerError)
		return
	}

	if tokens == nil {
		tokens = []string{}
	}

	resp, err := json.Marshal(struct {
		Tokens []string `json:"tokens"`
	}{
		Tokens: tokens,
	})
	if err != nil {
		panic(err)
	}

	w.Write(resp)
}

func (s *Server) handlePostToken(w http.ResponseWriter, ii *ImpersonateInfo, uid string) {
	tokens, err := s.listTokens(uid) // Check if the user has too many tokens
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to list tokens", http.StatusInternalServerError)
	}

	if len(tokens) >= *s.conf.TokenCountMax {
		http.Error(w, "Too many tokens", http.StatusForbidden)
		return
	}

	token := s.generateToken(uid)
	err = s.setToken(token, ii)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to set token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{\"token\":\"" + token + "\"}\n"))
}

func (s *Server) handleDeleteToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[len("/_/tokens"):]
	if token == "" || token == "/" {
		http.Error(w, "No token provided", http.StatusBadRequest)
		return
	}

	token = token[1:] // Remove leading slash

	err := s.deleteToken(token)
	if err != nil {
		http.Error(w, "Failed to delete token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{\"status\":\"success\"}\n"))
}
