PKG = github.com/k1LoW/gostyle
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: test gostyle

test:
	go test ./... -coverprofile=coverage.out -covermode=count

lint: gostyle
	golangci-lint run ./...
	govulncheck ./...

gostyle: build
	go vet -vettool=./gostyle -gostyle.config=$(PWD)/.gostyle.yml ./...

build:
	go build -ldflags="$(BUILD_LDFLAGS)" -o gostyle

depsdev:
	go install golang.org/x/vuln/cmd/govulncheck@latest

credits:
	go install github.com/Songmu/gocredits/cmd/gocredits@v1.0.0
	gocredits . > CREDITS
	cat _REFERENCE_STYLE_CREDITS >> CREDITS

prerelease_for_tagpr:
	$(MAKE) credits
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag
	goreleaser --clean

.PHONY: default test credits
