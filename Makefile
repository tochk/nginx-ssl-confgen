VERSION=1.0.0
DEB_DIR=tmp/nginx-confgen_${VERSION}_amd64
ARM_DEB_DIR=tmp/nginx-confgen_${VERSION}_amd64


generate:
	go generate ./...

build:
	go fmt ./...
	GOOS=linux go build -o bin/nginx-confgen cmd/nginx-confgen/main.go

build-arm:
	go fmt ./...
	GOOS=linux GOARCH=arm64 go build -o bin/nginx-confgen cmd/nginx-confgen/main.go

build-deb:
	rm -rf ${DEB_DIR}
	mkdir -p ${DEB_DIR}/usr/local/bin/
	cp bin/nginx-confgen ${DEB_DIR}/usr/local/bin/
	cp -r DEBIAN ${DEB_DIR}/
	dpkg-deb --build --root-owner-group ${DEB_DIR}

build-deb-arm:
	rm -rf ${ARM_DEB_DIR}
	mkdir -p ${ARM_DEB_DIR}/usr/local/bin/
	cp bin/nginx-confgen ${ARM_DEB_DIR}/usr/local/bin/
	cp -r DEBIAN ${ARM_DEB_DIR}/
	dpkg-deb --build --root-owner-group ${ARM_DEB_DIR}

all: generate build build-deb

all-arm: generate build-arm build-deb-arm