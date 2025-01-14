package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var anonymousImpersonateInfo = &ImpersonateInfo{
	Username: "system:anonymous",
	Group:    []string{"system:unauthenticated"},
}

func (s *Server) authenticate(req *http.Request) (*ImpersonateInfo, error) {
	token := req.Header.Get("Authorization")
	if token == "" {
		return anonymousImpersonateInfo, nil
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return nil, fmt.Errorf("invalid token")
	}

	token = token[7:]

	// SK token
	if strings.HasPrefix(token, "sk:") {
		return s.getToken(token)
	}

	// OAuth token
	oi, err := s.getOAuthUserInfo(token)
	if err != nil {
		return nil, err
	}
	if oi.ID == "" || oi.Name == "" {
		return nil, fmt.Errorf("invalid token")
	}
	return s.userInfoToImpersonateInfo(oi), nil
}

func (s *Server) handleWhoAmI(w http.ResponseWriter, r *http.Request) {
	ii, err := s.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := make(map[string]interface{})

	resp["username"] = ii.Username
	resp["group"] = ii.Group
	resp["extra"] = ii.Extra
	resp["uid"] = ii.UID

	respStr, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(respStr)
}
