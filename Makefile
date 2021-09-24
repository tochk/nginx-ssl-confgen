VERSION=1.0.0
DEB_DIR=tmp/nginx-confgen_${VERSION}_amd64

generate:
	go generate ./...

build:
	go fmt ./...
	GOOS=linux go build -o bin/nginx-confgen cmd/nginx-confgen/main.go

build-deb:
	rm -rf ${DEB_DIR}
	mkdir -p ${DEB_DIR}/usr/local/bin/
	cp bin/nginx-confgen ${DEB_DIR}/usr/local/bin/
	cp -r DEBIAN ${DEB_DIR}/
	dpkg-deb --build --root-owner-group ${DEB_DIR}

all: generate build build-deb