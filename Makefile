export GO15VENDOREXPERIMENT=1

DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | grep -v /vendor/)
NAME=kolide
WEBSITE=https://$(NAME).io
DESCRIPTION="Ask your environment questions"

BUILDVERSION=$(shell cat package/VERSION)
GO_VERSION=$(shell go version)

# Get the git commit
SHA=$(shell git rev-parse --short HEAD)
BUILD_COUNT=$(shell git rev-list --count HEAD)

BUILD_TAG="${BUILD_COUNT}.${SHA}"

build: banner lint generate
	@echo "Building $(NAME)..."
	@mkdir -p bin/
	@go build \
		-ldflags "-X main.build=${BUILD_TAG}" \
		${ARGS} \
		-o bin/$(NAME)

banner:
	@echo "$(NAME)"
	@echo "${GO_VERSION}"
	@echo "Go Path: ${GOPATH}"
	@echo

generate: cleanGoGenerate
	@echo "Running go generate..."
	@go generate $$(go list ./... | grep -v /vendor/)

deps:
	@echo "Installing build deps"
	go get -u github.com/golang/lint/golint
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

lint:
	@go vet  $$(go list ./... | grep -v /vendor/)
	@for pkg in $$(go list ./... |grep -v /vendor/ |grep -v /kuber/) ; do \
		golint -min_confidence=1 $$pkg ; \
		done

package: setup rpm64 deb64

setup: build strip
	@mkdir -p package/root/usr/bin/
	@mkdir -p package/root/etc/$(NAME)/
	@mkdir -p package/root/usr/lib/systemd/system/
	@mkdir -p package/output/
	@cp -R ./bin/$(NAME) package/root/usr/bin/$(NAME)
	@cp -R ./shared/$(NAME).toml package/root/etc/$(NAME)/$(NAME).toml
	@cp -R ./shared/$(NAME).service package/root/usr/lib/systemd/system/$(NAME).service
	@./bin/$(NAME) --version 2> package/VERSION

test:
	go list ./... | xargs -n1 go test

certs:
	mkdir -p tmp
	openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout tmp/$(NAME).key -out tmp/$(NAME).crt
	cp -R ./tmp/* /tmp/

certs-remote:
	mkdir -p tmp
	openssl req -x509 -sha256 -nodes -days 365 -subj "/C=us/ST=ks/L=kolide/O=kolide/CN=$$cn" -newkey rsa:2048 -keyout tmp/$(NAME).key -out tmp/$(NAME).crt
	cp -R ./tmp/* /tmp/

# docker dev
up:
	docker-compose up

down:
	docker-compose down

strip:
	strip bin/$(NAME)

rpm64: setup
	fpm -s dir -t rpm -n $(NAME) -v $(BUILDVERSION) -p package/output/ \
		--rpm-compression xz --rpm-os linux \
		--force \
		--before-install scripts/package/pre-inst.sh \
		--after-install scripts/package/post-inst.sh \
		--before-remove scripts/package/pre-rm.sh \
		--after-remove scripts/package/post-rm.sh \
		--url $(WEBSITE) \
		--description $(DESCRIPTION) \
		-m "kolide <engineering@kolide.co>" \
		--vendor "kolide" -a amd64 \
		--config-files etc/$(NAME)/$(NAME).toml \
		--exclude */**.gitkeep \
		package/root/=/

deb64: setup
	fpm -s dir -t deb -n $(NAME) -v $(BUILDVERSION) -p package/output/ \
		--force \
		--deb-compression xz \
		--before-install scripts/package/pre-inst.sh \
		--after-install scripts/package/post-inst.sh \
		--before-remove scripts/package/pre-rm.sh \
		--after-remove scripts/package/post-rm.sh \
		--url $(WEBSITE) \
		--description $(DESCRIPTION) \
		-m "kolide <engineering@kolide.co>" \
		--vendor "kolide" -a amd64 \
		--config-files etc/$(NAME)/$(NAME).toml \
		--exclude */**.gitkeep \
		package/root/=/

cleanGoGenerate:
	@rm -rf static/assets.go
	@rm -rf statis/bindata_assetfs.go

cleanServer: cleanGoGenerate
	@rm -rf doc/
	@rm -rf package/
	@rm -rf bin/
	@rm -rf tmp/

clean: cleanServer
	# @go clean -r -i

.PHONY: build test
