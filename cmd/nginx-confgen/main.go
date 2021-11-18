package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"nginx-confgen/internal/app/templates"
	"os"
	"os/exec"
	"strings"
)

var (
	nginxConfDir  = flag.String("nginx-conf-dir", "/etc/nginx/sites-available/", "nginx configuration directory")
	servers       = flag.String("servers", "", "comma-separeted domains list")
	proxyPass     = flag.String("proxy-pass", "", "proxy pass server")
	localDir      = flag.String("local-dir", "", "http files directory")
	sslFullChain  = flag.String("ssl-full-chain", "", "ssl full chain path")
	sslPrivateKey = flag.String("ssl-private-key", "", "ssl private key path")
	generateLeSSL = flag.Bool("generate-ssl", false, "generate letsencrypt certificate")
	email         = flag.String("email", "", "email for letsencrypt")
	agreeLeTos    = flag.Bool("agree-tos", false, "let's encrypt terms of service agreement")
)

func main() {
	flag.Parse()
	if *servers == "" {
		log.Fatal("servers flag must not be empty")
	}
	serversList := strings.Split(*servers, ",")
	for idx, srv := range serversList {
		serversList[idx] = strings.Trim(srv, " \t")
	}
	if *proxyPass == "" && *localDir == "" {
		log.Fatal("proxy pass or local dir not provided")
	}
	if *proxyPass != "" && *localDir != "" {
		log.Fatal("select one of the proxy pass or local dir")
	}
	if !*generateLeSSL {
		if *sslFullChain == "" || *sslPrivateKey == "" {
			log.Fatal("ssl full chain and private key must be not empty, or generate ssl flag enabled")
		}
	} else {
		if *email == "" {
			log.Fatal("you need to determine email for using letsencrypt")
		}
		if !*agreeLeTos {
			log.Fatal("you need to agree letsencrypt tos ")
		}
		*sslFullChain = fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", serversList[0])
		*sslPrivateKey = fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", serversList[0])
	}

	resultConfig := templates.NginxConfig{
		ProxyPass:         *proxyPass,
		ServerNames:       serversList,
		SSLFullChainPath:  *sslFullChain,
		SSLPrivateKeyPath: *sslPrivateKey,
		LocalDir:          *localDir,
	}

	if *generateLeSSL {
		if err := os.MkdirAll("/tmp/nginx-confgen/"+serversList[0], 0744); err != nil {
			log.Fatal("can't create temp directory:", err)
		}
		if err := os.WriteFile(*nginxConfDir+serversList[0]+".conf", []byte(templates.HttpConfig(resultConfig)), 0744); err != nil {
			log.Fatal("can't create http config file:", err)
		}
		if err := RestartNginx(); err != nil {
			log.Fatal("can't restart nginx:", err)
		}
		leArgs := make([]string, 0, (len(serversList)*2)+3)
		leArgs = append(leArgs, "--non-interactive", "--agree-tos", "--nginx", "--email", *email)
		for _, srv := range serversList {
			leArgs = append(leArgs, "-d", srv)
		}

		leCmd := exec.Command("certbot", leArgs...)
		if err := leCmd.Run(); err != nil {
			log.Fatal("can't generate cert with certbot:", err)
		}
	}

	if err := os.WriteFile(*nginxConfDir+"/"+serversList[0]+".conf", []byte(templates.HttpsConfig(resultConfig)), 0744); err != nil {
		log.Fatal("can't create https config file:", err)
	}

	if err := RestartNginx(); err != nil {
		log.Fatal("can't restart nginx:", err)
	}
}

func RestartNginx() error {
	nginxTestCmd := exec.Command("nginx", "-t")
	if err := nginxTestCmd.Run(); err != nil {
		return err
	}

	nginxRestartCmd := exec.Command("systemctl", "restart", "nginx")
	return nginxRestartCmd.Run()
}
