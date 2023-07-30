package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"nginx-ssl-confgen/internal/app/templates"
	"os"
	"os/exec"
	"strings"
)

var (
	nginxConfDir        = flag.String("nginx-conf-dir", "/etc/nginx/sites-available/", "nginx sites available directory")
	nginxConfDirEnabled = flag.String("nginx-conf-dir-enabled", "/etc/nginx/sites-enabled/", "nginx sites enabled directory")
	servers             = flag.String("servers", "", "comma-separeted domains list")
	proxyPass           = flag.String("proxy-pass", "", "proxy pass server")
	localDir            = flag.String("local-dir", "", "http files directory")
	sslFullChain        = flag.String("ssl-full-chain", "", "ssl full chain path")
	sslPrivateKey       = flag.String("ssl-private-key", "", "ssl private key path")
	generateLeSSL       = flag.Bool("generate-ssl", false, "generate letsencrypt certificate")
	email               = flag.String("email", "", "email for letsencrypt")
	agreeLeTos          = flag.Bool("agree-tos", false, "let's encrypt terms of service agreement")
)

func main() {
	flag.Parse()
	log.Infoln("checking provided flags")
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
			log.Fatal("you need to agree letsencrypt tos")
		}
		*sslFullChain = fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", serversList[0])
		*sslPrivateKey = fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", serversList[0])
	}

	log.Infoln("flags checked successfully")

	resultConfig := templates.NginxConfig{
		ProxyPass:         *proxyPass,
		ServerNames:       serversList,
		SSLFullChainPath:  *sslFullChain,
		SSLPrivateKeyPath: *sslPrivateKey,
		LocalDir:          *localDir,
	}

	if err := checkPrefix(resultConfig.ProxyPass); err != nil {
		log.Fatal(err)
	}

	if *generateLeSSL {
		log.Infoln("generating SSL cert")
		if err := os.MkdirAll("/tmp/nginx-ssl-confgen/"+serversList[0], 0755); err != nil {
			log.Fatal("can't create tmp directory:", err)
		}
		if err := os.WriteFile(*nginxConfDir+serversList[0]+".conf", []byte(templates.HttpConfig(resultConfig)), 0744); err != nil {
			log.Fatal("can't create http config file:", err)
		}
		if _, err := os.Stat(*nginxConfDirEnabled + serversList[0] + ".conf"); err == nil {
			err = os.Remove(*nginxConfDirEnabled + serversList[0] + ".conf")
			if err != nil {
				log.Fatal("can't remove symlink to http config file:", err)
			}
		}
		if err := os.Symlink(*nginxConfDir+serversList[0]+".conf", *nginxConfDirEnabled+serversList[0]+".conf"); err != nil {
			log.Fatal("can't create symlink to config:", err)
		}
		if err := RestartNginx(); err != nil {
			log.Fatal("can't restart nginx:", err)
		}
		leArgs := make([]string, 0, (len(serversList)*2)+3)
		leArgs = append(leArgs, "--non-interactive", "--agree-tos", "--nginx", "--email", *email)
		for _, srv := range serversList {
			leArgs = append(leArgs, "-d", srv)
		}

		log.Infoln("running certbot")
		leCmd := exec.Command("certbot", leArgs...)
		var leStdOut, leStdErr bytes.Buffer
		leCmd.Stdout = &leStdOut
		leCmd.Stderr = &leStdErr
		if err := leCmd.Run(); err != nil {
			log.Errorln("command was:", "certbot", leArgs)
			log.Errorln("certbot stdout:", leStdOut.String())
			log.Errorln("certbot stderr:", leStdErr.String())
			log.Fatal("can't generate cert with certbot:", err)
		}
		if err := os.Remove(*nginxConfDirEnabled + serversList[0] + ".conf"); err != nil {
			log.Warningln("can't remove", *nginxConfDirEnabled+serversList[0]+".conf", "file, skipping")
		}
		if err := os.Remove(*nginxConfDir + serversList[0] + ".conf"); err != nil {
			log.Warningln("can't remove", *nginxConfDir+serversList[0]+".conf", "file, skipping")
		}
		log.Infoln("certificates generated successfully")
	}

	if err := os.WriteFile(*nginxConfDir+serversList[0]+".conf", []byte(templates.HttpsConfig(resultConfig)), 0744); err != nil {
		log.Fatal("can't create https config file:", err)
	}

	if err := os.Symlink(*nginxConfDir+serversList[0]+".conf", *nginxConfDirEnabled+serversList[0]+".conf"); err != nil {
		log.Fatal("can't create symlink to config:", err)
	}

	if err := RestartNginx(); err != nil {
		log.Fatal("can't restart nginx:", err)
	}

	log.Infoln("completed!")
}

func RestartNginx() error {
	log.Infoln("checking nginx configuration")
	nginxTestCmd := exec.Command("nginx", "-t")
	if err := nginxTestCmd.Run(); err != nil {
		return err
	}

	log.Infoln("nginx config OK, restarting nginx")
	nginxRestartCmd := exec.Command("systemctl", "restart", "nginx")
	if err := nginxRestartCmd.Run(); err != nil {
		return err
	}

	log.Infoln("nginx restarted successfully")
	return nil
}

func checkPrefix(proxyPass string) error {
	if proxyPass == "" {
		return nil
	}

	if strings.Contains(proxyPass, "http://") ||
		strings.Contains(proxyPass, "https://") {
		return nil
	}

	return errors.New("proxy_pass must contain http:// or https:// prefix, got " + proxyPass)
}
