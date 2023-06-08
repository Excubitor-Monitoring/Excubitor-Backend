.PHONY: install-deps build run test test/coverage install package/deb components build-component

GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover
GOMOD=$(GO) mod
GOBUILD=$(GO) build
GORUN=$(GO) run

NPM=npm
NPMI=$(NPM) install
NPMBUILD=$(NPM) run build

EXCUBITOR_VERSION=0.0.1-alpha

install-deps:
	echo "Installing all go dependencies"
	$(GOMOD) download
components:
	git submodule init
	echo "Building frontend components"
	echo "Building CPU-Info component"
	make COMPDIR=components/CPU-Info MODNAME=cpu FILENAME=info.js build-component
	echo "Building CPU-Usage component"
	make COMPDIR=components/CPU-Usage MODNAME=cpu FILENAME=usage.js build-component
	echo "Building RAM-Usage component"
	make COMPDIR=components/RAM-Usage MODNAME=memory FILENAME=ram.js build-component
	echo "Building Swap-Usage component"
	make COMPDIR=components/Swap-Usage MODNAME=memory FILENAME=swap.js build-component
build-component:
	$(NPMI) --prefix $(COMPDIR)
	$(NPMBUILD) --prefix $(COMPDIR)
	mkdir -p internal/frontend/static/internal/$(MODNAME)
	mv $(COMPDIR)/dist/index.js internal/frontend/static/internal/$(MODNAME)/$(FILENAME)
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
	mkdir -p package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN
	cp package/deb/control package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN/control
	echo "Version: $(EXCUBITOR_VERSION)" >> package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN/control
	# Assemble package
	dpkg-deb --build --root-owner-group package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64
