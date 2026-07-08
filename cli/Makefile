BINARY  := scouti
DIST    := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

# Targets to cross-compile for `make release` (GOOS/GOARCH pairs).
PLATFORMS := \
	linux/amd64 linux/arm64 \
	darwin/amd64 darwin/arm64 \
	windows/amd64 windows/arm64

.PHONY: build release clean fmt vet

# Build a binary for the host platform into ./$(BINARY).
build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

# Cross-compile one static binary per platform into ./$(DIST).
release: clean
	@mkdir -p $(DIST)
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=""; \
		[ "$$os" = "windows" ] && ext=".exe"; \
		echo "  building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
			go build -trimpath -ldflags "$(LDFLAGS)" \
			-o $(DIST)/$(BINARY)-$$os-$$arch$$ext . || exit 1; \
	done
	@echo "-> $(DIST)/"

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(DIST) $(BINARY) $(BINARY).exe
