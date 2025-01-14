package main

import (
	"flag"
	"log"
	"time"

	"github.com/lcpu-club/kube-auth-proxy/internal/config"
	"github.com/lcpu-club/kube-auth-proxy/internal/server"
)

func main() {
	conf := &config.ServerConfig{}
	conf.Listen = flag.String("listen", ":8080", "Listen address")
	conf.Upstream = flag.String("upstream", "https://kubernetes.default.svc", "Upstream address")
	conf.Storage = flag.String("storage", "memory:", "Token storage type")
	conf.KubeSecretPath = flag.String("kube-secret-path", "/var/run/secrets/kubernetes.io/serviceaccount/token", "Path to the Kubernetes service account token")
	conf.OAuthAppID = flag.String("oauth-app-id", "", "OAuth App ID")
	conf.OAuthSecret = flag.String("oauth-secret", "", "OAuth App Secret")
	conf.OAuthCallback = flag.String("oauth-callback", "http://localhost:8080/_/oauth/callback", "OAuth Callback URL")
	conf.OAuthDefaultGroup = flag.String("oauth-default-group", "hpcgame:competitors", "Default group for OAuth users")
	conf.TokenExpiration = flag.Duration("token-expiration", 14*24*time.Hour, "Token expiration time")
	conf.TokenLength = flag.Int("token-length", 36, "Token length")
	conf.TokenCountMax = flag.Int("token-count-max", 128, "Maximum number of tokens per user")
	conf.TLSCertFile = flag.String("tls-cert-file", "", "TLS certificate file (empty to disable TLS)")
	conf.TLSKeyFile = flag.String("tls-key-file", "", "TLS key file (empty to disable TLS)")
	conf.KubeconfigTemplatePath = flag.String("kubeconfig-template-path", "kubeconfig.tmpl", "Path to the kubeconfig template file")

	flag.Parse()

	s := server.NewServer(conf)

	err := s.Init()
	if err != nil {
		log.Fatalln(err)
	}

	err = s.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
