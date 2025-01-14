package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lcpu-club/kube-auth-proxy/internal/utils"
	"golang.org/x/oauth2"
)

func (s *Server) oauthInit() {
	s.mux.HandleFunc("/_/oauth/callback", s.handleOAuthCallback)
	s.mux.HandleFunc("/_/oauth/redirect", s.handleOAuthRedirect)
	s.mux.HandleFunc("/_/oauth/userinfo", s.handleGetOAuthUserInfo)

	s.oauthConfig = &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://hpcgame.pku.edu.cn/oauth/authorize",
			TokenURL:  "https://hpcgame.pku.edu.cn/api/oauth/access_token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		ClientID:     *s.conf.OAuthAppID,
		ClientSecret: *s.conf.OAuthSecret,
		RedirectURL:  *s.conf.OAuthCallback,
	}
}

func (s *Server) handleOAuthRedirect(w http.ResponseWriter, r *http.Request) {
	// Redirect to the OAuth provider
	http.Redirect(w, r, s.oauthConfig.AuthCodeURL(
		utils.GenRandomStateString(),
	), http.StatusTemporaryRedirect)
}

func (s *Server) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	// TODO: FIXME: Check state string

	// Handle the OAuth callback
	code := r.FormValue("code")
	token, err := s.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		http.Redirect(w, r, "redirect", http.StatusTemporaryRedirect)
		w.Write([]byte("Failed to exchange token"))
		return
	}

	// TODO: Write reconcile User object logic

	// Redirect to the UI
	http.Redirect(w, r, fmt.Sprintf("../ui/#/auth/token/%s", token.AccessToken), http.StatusTemporaryRedirect)
	w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token.AccessToken)))
}

func (s *Server) getOAuthUserInfo(token string) (*OAuthUserInfo, error) {
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	uid, err := utils.ExtractUIDFromJWT(token)
	if err != nil {
		return nil, err
	}

	userInfoURL := fmt.Sprintf("https://hpcgame.pku.edu.cn/api/user/%s/profile", uid)
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "\"error\":\"Unauthorized\"") {
		return nil, fmt.Errorf("Unauthorized")
	}

	rslt := &OAuthUserInfo{}
	err = json.Unmarshal(body, rslt)
	if err != nil {
		return nil, err
	}
	rslt.ID = uid

	return rslt, nil
}

type OAuthUserInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Realname     string   `json:"realname"`
	Telephone    string   `json:"telephone"`
	School       string   `json:"school"`
	StudentGrade string   `json:"studentGrade"`
	Verified     []string `json:"verified"`
}

func (s *Server) handleGetOAuthUserInfo(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	userInfo, err := s.getOAuthUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, "Failed to marshal user info", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (s *Server) userInfoToImpersonateInfo(userInfo *OAuthUserInfo) *ImpersonateInfo {
	ii := &ImpersonateInfo{
		UID: userInfo.ID,
		// TODO: More elegant way to handle ID
		Username: userInfo.ID,
		Group:    []string{*s.conf.OAuthDefaultGroup},
		Extra: map[string]string{
			"name":         userInfo.Name,
			"email":        userInfo.Email,
			"realname":     userInfo.Realname,
			"telephone":    userInfo.Telephone,
			"school":       userInfo.School,
			"studentGrade": userInfo.StudentGrade,
		},
	}
	return ii
}
