# Nginx config generator (nginx-ssl-confgen)

Simple nginx config generator with let's encrypt support

## Prerequisites

- nginx
- certbot (if you want to use let's encrypt certificate generation)
- python3-certbot-nginx (if you want to use let's encrypt certificate generation)

### Ubuntu server installation

```shell
apt-get update
apt-get install nginx certbot python3-certbot-nginx
```

## Installation

Go to releases page and download `deb` package.

## Usage

```
Usage of nginx-ssl-confgen:
  -agree-tos
        let's encrypt terms of service agreement
  -email string
        email for letsencrypt
  -generate-ssl
        generate letsencrypt certificate
  -local-dir string
        http files directory
  -nginx-conf-dir string
        nginx sites available directory (default "/etc/nginx/sites-available/")
  -nginx-conf-dir-enabled string
        nginx sites enabled directory (default "/etc/nginx/sites-enabled/")
  -proxy-pass string
        proxy pass server
  -servers string
        comma-separeted domains list
  -ssl-full-chain string
        ssl full chain path
  -ssl-private-key string
        ssl private key path
```

Example (with let's encrypt certificate and proxy pass):
```shell
nginx-ssl-confgen -servers=example.com -proxy-pass=http://localhost:8080 -generate-ssl -email=me@example.com -agree-tos
```

### Crontab setup (for certbot renewal)

If your certbot package does not have `certbot.timer`, add the following line to the crontab for automatic certificate renewal:

```cronexp
0 6 * * * certbot renew
```

## Building

### Prerequisites

- go 1.20

### Building binary

```shell
make build-amd64 # amd64 binary
```
OR
```shell
make build-arm64 # arm64 binary
```

### Building deb package

```shell
make build-deb-amd64 # amd64 deb package
```
OR
```shell
make build-deb-arm64 # arm64 deb package
```

# TODO

- [ ] Tests
- [ ] Linters
- [ ] RPM builds