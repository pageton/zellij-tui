# zellij-tui build tasks
# Requires: go, just (https://github.com/casey/just)

binary := "zellij-tui"
goflags := "-trimpath"
version := `git describe --tags --always --dirty 2>/dev/null || echo dev`
ldflags := "-s -w -X main.version=" + version

# default: build optimized binary
default: build

# ── build ──────────────────────────────────────────────

# build with CGO disabled (static, fast, version stamped)
build:
    CGO_ENABLED=0 GOAMD64=v3 go build {{goflags}} -ldflags="{{ldflags}}" -o {{binary}} .

# build with debug symbols (for profiling / delve)
debug:
    go build -gcflags="all=-N -l" -o {{binary}} .

# cross-compile for all supported targets
build-all: (build-os-arch "linux" "amd64") (build-os-arch "linux" "arm64") (build-os-arch "darwin" "arm64") (build-os-arch "darwin" "amd64")

# ── test ───────────────────────────────────────────────

# run all tests
test:
    go test ./... -count=1

# run tests with verbose output
test-v:
    go test ./... -count=1 -v

# run tests with race detector
test-race:
    go test ./... -count=1 -race

# run only the zellij package tests (fast path for active dev)
test-fast:
    go test ./internal/zellij/ -count=1 -v

# ── formatting ─────────────────────────────────────────

# format all Go files (gofmt + goimports)
fmt:
    gofmt -w .
    goimports -w .

# check formatting without modifying files
fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        echo "Files need gofmt formatting:" >&2
        echo "$unformatted" >&2
        exit 1
    fi
    unformatted=$(goimports -l .)
    if [ -n "$unformatted" ]; then
        echo "Files need goimports formatting:" >&2
        echo "$unformatted" >&2
        exit 1
    fi

# ── lint & static analysis ────────────────────────────

# run go vet
vet:
    go vet ./...

# run staticcheck (requires staticcheck)
staticcheck:
    staticcheck ./...

# run golangci-lint with project config (requires golangci-lint)
lint:
    golangci-lint run --timeout=3m

# ── combined checks ───────────────────────────────────

# full validation gate (format check + lint + vet + test)
check: fmt-check lint vet test

# mirror the CI pipeline locally (format + lint + vet + race test + build)
ci: fmt-check lint vet test-race build

# ── module maintenance ────────────────────────────────

# tidy and verify module graph
tidy:
    go mod tidy
    go mod verify

# ── clean ──────────────────────────────────────────────

# clean build artifacts
clean:
    rm -f {{binary}}

# clean go build cache
clean-cache: clean
    go clean -cache

# ── utilities ─────────────────────────────────────────

# benchmark: cold build with clean cache
bench-cold: clean-cache
    @echo "=== Cold build ==="
    @time sh -c 'CGO_ENABLED=0 go build {{goflags}} -ldflags="{{ldflags}}" -o {{binary}} .'

# show binary size and version
size: build
    @ls -lh {{binary}}
    @./{{binary}} --version

# install to GOPATH/bin
install: build
    cp {{binary}} "$(go env GOPATH)/bin/"
    @echo "installed to $(go env GOPATH)/bin/{{binary}}"

# ── private recipes ───────────────────────────────────

[private]
build-os-arch goos goarch:
    #!/usr/bin/env bash
    set -euo pipefail
    out="{{binary}}-{{goos}}-{{goarch}}"
    echo "Building $out..."
    CGO_ENABLED=0 GOOS={{goos}} GOARCH={{goarch}} \
        go build {{goflags}} -ldflags="{{ldflags}}" -o "$out" .
