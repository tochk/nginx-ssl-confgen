VERSION=1.1.0
AMD_DEB_DIR=tmp/nginx-ssl-confgen_${VERSION}_amd64
ARM_DEB_DIR=tmp/nginx-ssl-confgen_${VERSION}_arm64


generate:
	go generate ./...

build-amd64:
	go fmt ./...
	GOOS=linux GOARCH=amd64 go build -o bin/nginx-ssl-confgen-amd64 cmd/nginx-ssl-confgen/main.go

build-arm64:
	go fmt ./...
	GOOS=linux GOARCH=arm64 go build -o bin/nginx-ssl-confgen-arm64 cmd/nginx-ssl-confgen/main.go

build-deb-arm64: build-amd64
	rm -rf ${ARM_DEB_DIR}
	mkdir -p ${ARM_DEB_DIR}/usr/local/bin/
	cp bin/nginx-ssl-confgen-arm64 ${ARM_DEB_DIR}/usr/local/bin/nginx-ssl-confgen
	cp -r deb/DEBIAN-arm64 ${ARM_DEB_DIR}/DEBIAN
	dpkg-deb --build --root-owner-group ${ARM_DEB_DIR}

build-deb-amd64: build-amd64
	rm -rf ${AMD_DEB_DIR}
	mkdir -p ${AMD_DEB_DIR}/usr/local/bin/
	cp bin/nginx-ssl-confgen-amd64 ${AMD_DEB_DIR}/usr/local/bin/nginx-ssl-confgen
	cp -r deb/DEBIAN-amd64 ${AMD_DEB_DIR}/DEBIAN
	dpkg-deb --build --root-owner-group ${AMD_DEB_DIR}

all-amd: generate build-deb-amd64

all-arm: generate build-deb-arm64

all: all-amd all-arm