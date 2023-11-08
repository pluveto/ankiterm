PROJECT_NAME=ankiterm-app
BIN_NAMES=ankiterm
APP=ankiterm
GOARCHS=amd64 386 arm arm64
GOARCHS_MAC=amd64 arm64
INSTALL_DIR=/usr/local/bin
LD_FLAGS=-ldflags="-X 'main.Version=$(shell git describe --tags --long 2>/dev/null || echo unknown)'"

dev: linux

default: all

all: windows linux mac

prepare:
	@mkdir -p dist

windows: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS); do \
			mkdir -p dist/windows_$$GOARCH; \
			GOOS=windows GOARCH=$$GOARCH go build $(LD_FLAGS) -o dist/windows_$$GOARCH/$$BIN_NAME.exe cmd/$$BIN_NAME/main.go; \
		done \
	done

linux: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS); do \
			mkdir -p dist/linux_$$GOARCH; \
			GOOS=linux GOARCH=$$GOARCH go build $(LD_FLAGS) -o dist/linux_$$GOARCH/$$BIN_NAME cmd/$$BIN_NAME/main.go; \
		done \
	done

mac: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS_MAC); do \
			mkdir -p dist/mac_$$GOARCH; \
			GOOS=darwin GOARCH=$$GOARCH go build $(LD_FLAGS) -o dist/mac_$$GOARCH/$$BIN_NAME cmd/$$BIN_NAME/main.go; \
		done \
	done

package: all
	for GOARCH in $(GOARCHS); do \
		zip -q -r dist/$(PROJECT_NAME)-windows-$$GOARCH.zip dist/windows_$$GOARCH/; \
		zip -q -r dist/$(PROJECT_NAME)-linux-$$GOARCH.zip dist/linux_$$GOARCH/; \
	done

	for GOARCH in $(GOARCHS_MAC); do \
		zip -q -r dist/$(PROJECT_NAME)-mac-$$GOARCH.zip dist/mac_$$GOARCH/; \
	done

	ARCH_RELEASE_DIRS=$$(find dist -type d -name "*_*"); \
	for ARCH_RELEASE_DIR in $$ARCH_RELEASE_DIRS; do \
		cp conf/config.default.toml $$ARCH_RELEASE_DIR/config.toml; \
		rm -rfd $$ARCH_RELEASE_DIR; \
	done

test:
	go test -v ./...

run:
	go run cmd/$(APP)/main.go $(ARGS)

clean:
	rm -rfd dist

install:
	go build $(LD_FLAGS) -o dist/$(APP)-install cmd/$(APP)/main.go
	sudo mv dist/$(APP)-install $(INSTALL_DIR)/$(APP)
	sudo chmod +x $(INSTALL_DIR)/$(APP)

uninstall:
	rm -rfd $(INSTALL_DIR)/$(APP)

.PHONY: all, default, clean
