package config

import "time"

type ServerConfig struct {
	Listen     *string
	Upstream   *string
	Storage    *string
	UIDistPath *string

	KubeSecretPath *string

	OAuthCallback *string
	OAuthAppID    *string
	OAuthSecret   *string

	OAuthDefaultGroup *string

	TokenExpiration *time.Duration
	TokenLength     *int
	TokenCountMax   *int

	TLSCertFile *string
	TLSKeyFile  *string

	KubeconfigTemplatePath *string
}
