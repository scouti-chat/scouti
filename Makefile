# Scouti devkit — release packaging.
#
# `make` (a.k.a. `make dist`) builds everything we ship into ./dist at the repo
# root, ready to grab and upload to a GitHub release by hand:
#   - scouti-<os>-<arch>[.exe]      raw CLI binary per platform (direct download)
#   - scouti-<os>-<arch>.tar.gz     CLI archive per platform (holds `scouti[.exe]`)
#   - skill.tar.gz                  the installable skill bundle
#
# The Go module lives in ./cli; the skill sources live in ./skill (run
# `sc sync_assets` in the main repo first so skill/guide.md + skill/api.md exist).

BINARY    := scouti
DIST      := dist
CLI_DIR   := cli
SKILL_DIR := skill
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS   := -s -w -X main.version=$(VERSION)

# GOOS/GOARCH pairs to cross-compile.
PLATFORMS := \
	linux/amd64 linux/arm64 \
	darwin/amd64 darwin/arm64 \
	windows/amd64 windows/arm64

# Files bundled into skill.tar.gz (authored SKILL.md + docs synced from the main repo).
SKILL_FILES := SKILL.md guide.md api.md

# Serialize so `clean` never races the build targets under `make -j`.
.NOTPARALLEL:
.PHONY: dist cli skill build fmt vet clean

# Default: build the full, upload-ready release set into ./dist.
dist: clean cli skill
	@echo ""
	@echo "-> $(DIST)/"
	@ls -1 $(DIST)

# One raw binary + one .tar.gz (holding `scouti[.exe]`) per platform.
cli:
	@mkdir -p $(DIST)
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=""; \
		[ "$$os" = "windows" ] && ext=".exe"; \
		echo "  cli    $$os/$$arch"; \
		stage=$$(mktemp -d); \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
			go build -C $(CLI_DIR) -trimpath -ldflags "$(LDFLAGS)" \
			-o $$stage/$(BINARY)$$ext . || exit 1; \
		cp $$stage/$(BINARY)$$ext $(DIST)/$(BINARY)-$$os-$$arch$$ext; \
		tar -czf $(DIST)/$(BINARY)-$$os-$$arch.tar.gz -C $$stage $(BINARY)$$ext; \
		rm -rf $$stage; \
	done

# Self-contained skill bundle. Extracted flat into an agent's skills dir.
skill:
	@mkdir -p $(DIST)
	@echo "  skill  skill.tar.gz"
	@tar -czf $(DIST)/skill.tar.gz -C $(SKILL_DIR) $(SKILL_FILES)

# Host build for local use (outputs ./cli/scouti; not part of a release).
build:
	cd $(CLI_DIR) && go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) .

fmt:
	cd $(CLI_DIR) && go fmt ./...

vet:
	cd $(CLI_DIR) && go vet ./...

clean:
	rm -rf $(DIST)
