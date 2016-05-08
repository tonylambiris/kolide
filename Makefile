export GO15VENDOREXPERIMENT=1

DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | grep -v /vendor/)
NAME="kolide"
WEBSITE="https://${NAME}.io"
DESCRIPTION="Ask your environment questions"

BUILDVERSION=$(shell cat VERSION)
GO_VERSION=$(shell go version)

# Get the git commit
SHA=$(shell git rev-parse --short HEAD)
BUILD_COUNT=$(shell git rev-list --count HEAD)

BUILD_TAG="${BUILD_COUNT}.${SHA}"

build: banner lint generate
	@echo "Building ${NAME}..."
	@mkdir -p bin/
	@go build \
    -ldflags "-X main.build=${BUILD_TAG}" \
		${ARGS} \
    -o bin/${NAME}

banner:
	@echo "${NAME}"
	@echo "${GO_VERSION}"
	@echo "Go Path: ${GOPATH}"
	@echo

generate: cleanGoGenerate
	@echo "Running go generate..."
	@go generate $$(go list ./... | grep -v /vendor/)

lint:
	@go vet  $$(go list ./... | grep -v /vendor/)
	@for pkg in $$(go list ./... |grep -v /vendor/ |grep -v /kuber/) ; do \
		golint -min_confidence=1 $$pkg ; \
		done

client: clientDeps
	@echo "Building client..."
	@cd client && gulp prod

package: setup strip rpm64

setup:
	@mkdir -p package/root/opt/${NAME}/bin/
	@mkdir -p package/root/etc/${NAME}/
	@mkdir -p package/output/
	@cp -R ./bin/${NAME} package/root/opt/${NAME}/bin
	@cp -R ./shared/${NAME}.toml package/root/etc/${NAME}/${NAME}.toml
	@./bin/${NAME} --version 2> VERSION

test:
	go list ./... | xargs -n1 go test

certs:
	mkdir -p tmp
	openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout tmp/${NAME}.key -out tmp/${NAME}.crt

cert-remote:
	mkdir -p tmp
	openssl req -x509 -sha256 -nodes -days 365 -subj "/C=us/ST=ks/L=kolide/O=kolide/CN=209.6.37.74" -newkey rsa:2048 -keyout tmp/${NAME}.key -out tmp/${NAME}.crt 

# docker dev
up:
	docker-compose up

down:
	docker-compose down

strip:
	strip bin/${NAME}

rpm64:
	fpm -s dir -t rpm -n $(NAME) -v $(BUILDVERSION) -p package/output/${NAME}-$(BUILDVERSION)-amd64.rpm \
		--rpm-compression bzip2 --rpm-os linux \
		--rpm-user ${NAME} --rpm-group ${NAME} \
		--force \
		--after-install scripts/rpm/post-install.sh \
		--before-remove scripts/rpm/pre-rm.sh \
		--after-remove scripts/rpm/post-rm.sh \
		--url $(WEBSITE) \
		--description $(DESCRIPTION) \
		-m "Dustin Wills Webber <dustin.webber@gmail.com>" \
		--vendor "Dustin Willis Webber" -a amd64 \
		--config-files etc/${NAME}/${NAME}.toml \
		--exclude */**.gitkeep \
		package/root/=/

cleanGoGenerate:
	@rm -rf static/assets.go
	@rm -rf statis/bindata_assetfs.go

cleanServer: cleanGoGenerate
	@rm -rf doc/
	@rm -rf package/
	@rm -rf bin/
	@rm -rf VERSION
	@rm -rf tmp/

clean: cleanServer
	# @go clean -r -i

.PHONY: build test
