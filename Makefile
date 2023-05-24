.PHONY: install-deps build run test test/coverage package/deb

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
install:
	make build
	# Install Configuration file
	mkdir -p $(DESTDIR)/etc/excubitor
	install -m 0755 config.sample.yml $(DESTDIR)/etc/excubitor/config.yml
	# Install binary
	mkdir -p $(DESTDIR)/opt/excubitor/bin
	install -m 0755 bin/excubitor-backend $(DESTDIR)/opt/excubitor/bin/excubitor-backend
	# Install systemd unit file
	mkdir -p $(DESTDIR)/etc/systemd/system
	install -m 0755 package/systemd/excubitor.service $(DESTDIR)/etc/systemd/system/excubitor.service
package/deb:
	make DESTDIR=package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/ install
	# Copying control file and adding version
	cp package/deb/control package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN
	echo "Version: $(EXCUBITOR_VERSION)" >> package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN/control
	# Assemble package
	dpkg-deb --build --root-owner-group package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64
