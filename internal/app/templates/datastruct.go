package templates

type NginxConfig struct {
	ProxyPass         string
	ServerNames       []string
	SSLFullChainPath  string
	SSLPrivateKeyPath string
	LocalDir          string
}
