.PHONY: install-deps build run test test/coverage install package/deb components build-component rebuild

GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover
GOMOD=$(GO) mod
GOBUILD=$(GO) build -ldflags="-s -w"
GORUN=$(GO) run

NPM=yarn
NPMI=install
NPMBUILD=run build

UPX=upx

EXCUBITOR_VERSION=0.0.1-alpha

install-deps:
	echo "Installing all go dependencies"
	$(GOMOD) download
components:
	git submodule init
	@echo "Building frontend components"
	@echo "Building CPU-Info component"
	make components/CPU-Info/dist
	@echo "Building CPU-Clock-History component"
	make components/CPU-Clock-History/dist
	@echo "Building CPU-Usage component"
	make components/CPU-Usage/dist
	@echo "Building CPU-Usage-History component"
	make components/CPU-Usage-History/dist
	@echo "Building RAM-Usage component"
	make components/RAM-Usage/dist
	@echo "Building RAM-Usage-History component"
	make components/RAM-Usage-History/dist
	@echo "Building Swap-Usage component"
	make components/Swap-Usage/dist
	@echo "Building Swap-Usage-History component"
	make components/Swap-Usage-History/dist
components/CPU-Info/dist:
	make COMPDIR=components/CPU-Info MODNAME=cpu FILENAME=info.js build-component
components/CPU-Clock-History/dist:
	make COMPDIR=components/CPU-Clock-History MODNAME=cpu FILENAME=clock-history.js build-component
components/CPU-Usage/dist:
	make COMPDIR=components/CPU-Usage MODNAME=cpu FILENAME=usage.js build-component
components/CPU-Usage-History/dist:
	make COMPDIR=components/CPU-Usage-History MODNAME=cpu FILENAME=usage-history.js build-component
components/RAM-Usage/dist:
	make COMPDIR=components/RAM-Usage MODNAME=memory FILENAME=ram.js build-component
components/RAM-Usage-History/dist:
	make COMPDIR=components/RAM-Usage-History MODNAME=memory FILENAME=ram-history.js build-component
components/Swap-Usage/dist:
	make COMPDIR=components/Swap-Usage MODNAME=memory FILENAME=swap.js build-component
components/Swap-Usage-History/dist:
	make COMPDIR=components/Swap-Usage-History MODNAME=memory FILENAME=swap-history.js build-component
build-component:
	$(NPM) --cwd $(COMPDIR) $(NPMI)
	$(NPM) --cwd $(COMPDIR) $(NPMBUILD)
	mkdir -p internal/frontend/static/internal/$(MODNAME)
	mv $(COMPDIR)/dist/index.js internal/frontend/static/internal/$(MODNAME)/$(FILENAME)
build:
	make components
	@echo "Compiling project for current platform"
	$(GOBUILD) -o bin/excubitor-backend ./cmd/main.go
	@if [ "$(USE_UPX)" = "true" ]; then echo "Using UPX to compress the binary..."; $(UPX) bin/excubitor-backend; fi
rebuild:
	make clean
	make build
clean:
	@echo "Removing binary packages"
	rm -rf bin/excubitor-backend
	@echo "Removing built javascript files"
	rm -rf components/*/dist
	rm -rf internal/frontend/static/internal/*
	@echo "Removing javascript dependencies"
	rm -rf components/*/node_modules
	@echo "Removing packaging files"
	rm -rf package/deb/excubitor_*
	@echo "Removing coverage reports"
	rm -rf coverage.out
run:
	make components
	$(GORUN) cmd/main.go
test:
	make components
	$(GOTEST) -v ./...
test/coverage:
	make components
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out
install:
	make build
	# Install Configuration file
	mkdir -p $(DESTDIR)/etc/excubitor
	install -m 0644 config.sample.yml $(DESTDIR)/etc/excubitor/config.yml
	# Install binary
	mkdir -p $(DESTDIR)/opt/excubitor/bin
	install -m 0755 bin/excubitor-backend $(DESTDIR)/opt/excubitor/bin/excubitor-backend
	# Install systemd unit file
	mkdir -p $(DESTDIR)/lib/systemd/system
	install -m 0644 package/systemd/excubitor.service $(DESTDIR)/lib/systemd/system/excubitor.service
	# Install copyright file
	mkdir -p $(DESTDIR)/usr/share/doc/excubitor
	install -m 0644 package/deb/copyright $(DESTDIR)/usr/share/doc/excubitor
package/deb:
	make DESTDIR=package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64 install
	# Copying control file and adding version
	mkdir -p package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN
	cp package/deb/control package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN/control
	cp package/deb/conffiles package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN/conffiles
	@echo "Version: $(EXCUBITOR_VERSION)" >> package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64/DEBIAN/control
	# Assemble package
	dpkg-deb --build --root-owner-group package/deb/excubitor_$(EXCUBITOR_VERSION)_amd64
