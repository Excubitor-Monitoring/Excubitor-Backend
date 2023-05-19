GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover
GOMOD=$(GO) mod
GOBUILD=$(GO) build
GORUN=$(GO) run

EXCUBITOR_VERSION=0.0.1-alpha

install-deps:
	echo "Installing all go dependencies"
	$(GOMOD) download
build:
	echo "Compiling project for current platform"
	$(GOBUILD) -o bin/excubitor-backend ./cmd/main.go
run:
	$(GORUN) cmd/main.go
test:
	$(GOTEST) -v ./...
test/coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out
package/deb:
	make build
	# Add binary to package
	mkdir -p build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/opt/excubitor/bin
	cp bin/excubitor-backend build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/opt/excubitor/bin
	# Add systemd unit file to package
	mkdir -p build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/etc/systemd/system
	cp build/systemd/excubitor.service build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/etc/systemd/system/
	# Add config file to package
	mkdir -p build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/etc/excubitor
	cp config.sample.yml build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/etc/excubitor
	# Assemble package
	dpkg-deb --build --root-owner-group build/deb/excubitor_$(EXCUBITOR_VERSION)_amd64
