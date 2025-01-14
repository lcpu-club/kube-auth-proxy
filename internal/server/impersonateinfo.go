package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ImpersonateInfo struct {
	UID      string            `json:"i"`
	Username string            `json:"u"`
	Group    []string          `json:"g"`
	Extra    map[string]string `json:"e"`
}

func ImpersonateInfoFromString(s string) (*ImpersonateInfo, error) {
	ii := &ImpersonateInfo{}
	err := json.Unmarshal([]byte(s), ii)
	if err != nil {
		return nil, err
	}
	return ii, nil
}

func (ii *ImpersonateInfo) String() string {
	s, err := json.Marshal(ii)
	if err != nil {
		panic(err) // Should never fail
	}
	return string(s)
}

func (ii *ImpersonateInfo) Clean(req *http.Request) {
	for k := range req.Header {
		if strings.HasPrefix(k, "Impersonate-") {
			req.Header.Del(k)
		}
	}
}

func (ii *ImpersonateInfo) Render(req *http.Request) {
	req.Header.Set("Impersonate-User", ii.Username)
	req.Header.Set("Impersonate-Group", ii.Group[0])
	req.Header.Set("Impersonate-Uid", ii.UID)
	for k, v := range ii.Extra {
		req.Header.Set("Impersonate-Extra-"+k, v)
	}
}
