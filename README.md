# Nginx config generator (nginx-ssl-confgen)

Simple nginx config generator with let's encrypt support

## Prerequisites

- Go 1.16+
- Nginx
- You need to have certbot installed, if you want to use let's encrypt certificate generation.

## Installation

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
nginx-ssl-confgen -servers=tochk.net -proxy-pass=http://localhost:8080 -generate-ssl -email=me@tochk.net -agree-tos
```

# TODO

- [ ] Tests
- [ ] Linters
- [ ] RPM builds