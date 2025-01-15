package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"
	"time"

	"github.com/lcpu-club/kube-auth-proxy/internal/config"
	"github.com/lcpu-club/kube-auth-proxy/internal/utils"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Server struct {
	mux      *http.ServeMux
	conf     *config.ServerConfig
	upstream *url.URL
	rev      *httputil.ReverseProxy

	sm   *utils.SecretManager
	stor TokenStorage

	oauthConfig *oauth2.Config

	kubeconfigTemplate *template.Template

	kubeClient *dynamic.DynamicClient
	scheme     *runtime.Scheme
	userGVR    schema.GroupVersionResource
	userGVK    schema.GroupVersionKind
}

func NewServer(conf *config.ServerConfig) *Server {
	return &Server{
		mux:  http.NewServeMux(),
		conf: conf,
		sm:   utils.NewSecretManager(*conf.KubeSecretPath),
	}
}

func (s *Server) Init() (err error) {
	s.stor, err = NewTokenStorage(*s.conf.Storage)
	if err != nil {
		return err
	}

	s.upstream, err = url.Parse(*s.conf.Upstream)
	if err != nil {
		return err
	}

	err = s.initKube()
	if err != nil {
		return err
	}

	s.rev = httputil.NewSingleHostReverseProxy(s.upstream)
	s.rev.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs: s.sm.GetRootCAs(),
		},
	}

	s.oauthInit()
	s.initToken()
	s.mux.Handle("/_/whoami", http.HandlerFunc(s.handleWhoAmI))

	s.mux.Handle("/_/ui/", http.StripPrefix("/_/ui/", http.FileServer(http.Dir(*s.conf.UIDistPath))))
	s.mux.HandleFunc("/", s.HandleProxy)

	return nil
}

func (s *Server) HandleProxy(w http.ResponseWriter, r *http.Request) {
	ii, err := s.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ii.Clean(r)
	ii.Render(r)

	r.Header.Set("Authorization", "Bearer "+s.sm.GetToken())
	s.rev.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	log.Println("Listening on", *s.conf.Listen)
	if *s.conf.TLSCertFile != "" && *s.conf.TLSKeyFile != "" {
		return http.ListenAndServeTLS(*s.conf.Listen, *s.conf.TLSCertFile, *s.conf.TLSKeyFile, s.mux)
	}
	return http.ListenAndServe(*s.conf.Listen, s.mux)
}
