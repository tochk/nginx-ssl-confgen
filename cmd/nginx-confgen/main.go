package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

var (
	nginxConfDir  = flag.String("nginx-conf-dir", "/etc/nginx/", "nginx configuration directory")
	servers       = flag.String("servers", "", "comma-separeted domains list")
	proxyPass     = flag.String("proxy-pass", "", "proxy pass server")
	localDir      = flag.String("local-dir", "", "http files directory")
	sslFullChain  = flag.String("ssl-full-chain", "", "ssl full chain path")
	sslPrivateKey = flag.String("ssl-private-key", "", "ssl private key path")
)

func main() {
	if *servers == "" {
		log.Fatal("servers flag must not be empty")
	}
	if *proxyPass == "" && *localDir == "" {
		log.Fatal("proxy pass or local dir not provided")
	}
	if *proxyPass != "" && *localDir != "" {
		log.Fatal("select one of the proxy pass or local dir")
	}

}
