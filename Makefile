BINARY := dist/wht-ls-github
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
RELEASE_GOOS := darwin
RELEASE_GOARCH := arm64
PKG := ./src

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

.PHONY: all build test clean help version show-version release release-minor release-major computed-version do-release ci-release

all: build

help:
	@echo '使い方'
	@echo '  make / make build     このマシン向けに $(BINARY) をビルドする'
	@echo '  make test             テストを実行する'
	@echo '  make clean            dist/ を消す'
	@echo '  make show-version     現在のタグと次の patch / minor / major を表示する'
	@echo '  make version          show-version と同じ'
	@echo '  make release          パッチを上げてタグ・ビルド・GitHub Releases する'
	@echo '  make release-minor    マイナーを上げてリリースする'
	@echo '  make release-major    メジャーを上げてリリースする'
	@echo '  make ci-release       Actions 用。KIND=patch|minor|major で do-release する'
	@echo
	@echo 'make build の版数は git describe（無ければ dev）。明示するとき: make build VERSION=v0.1.0'
	@echo 'リリース成果物は $(RELEASE_GOOS)/$(RELEASE_GOARCH) 固定。'

build:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(PKG)

test:
	go test $(PKG)/...

clean:
	rm -rf dist

version: show-version

show-version:
	@git fetch origin --tags >/dev/null 2>&1 || true
	@last=$$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n1); \
	if [ -z "$$last" ]; then \
		echo '現在: （タグなし）'; \
		echo '次の patch: v0.1.0'; \
		echo '次の minor: v0.1.0'; \
		echo '次の major: v0.1.0'; \
		exit 0; \
	fi; \
	ver=$${last#v}; \
	major=$$(printf '%s\n' "$$ver" | cut -d. -f1); \
	minor=$$(printf '%s\n' "$$ver" | cut -d. -f2); \
	patch=$$(printf '%s\n' "$$ver" | cut -d. -f3); \
	echo "現在: $$last"; \
	echo "次の patch: v$$major.$$minor.$$((patch + 1))"; \
	echo "次の minor: v$$major.$$((minor + 1)).0"; \
	echo "次の major: v$$((major + 1)).0.0"

computed-version:
	@git fetch origin --tags >/dev/null 2>&1 || true
	@last=$$(git tag -l 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | head -n1); \
	kind="$(KIND)"; \
	if [ -z "$$last" ]; then echo v0.1.0; exit 0; fi; \
	ver=$${last#v}; \
	major=$$(printf '%s\n' "$$ver" | cut -d. -f1); \
	minor=$$(printf '%s\n' "$$ver" | cut -d. -f2); \
	patch=$$(printf '%s\n' "$$ver" | cut -d. -f3); \
	case "$$kind" in \
		major) printf 'v%s.0.0\n' $$((major + 1)) ;; \
		minor) printf 'v%s.%s.0\n' "$$major" $$((minor + 1)) ;; \
		patch) printf 'v%s.%s.%s\n' "$$major" "$$minor" $$((patch + 1)) ;; \
		*) echo '内部エラー: KIND が不正です。' >&2; exit 1 ;; \
	esac

release:
	@$(MAKE) do-release KIND=patch

release-minor:
	@$(MAKE) do-release KIND=minor

release-major:
	@$(MAKE) do-release KIND=major

ci-release:
	@$(MAKE) do-release KIND="$(KIND)"

do-release:
	@set -e; \
	test -z "$$(git status --porcelain)" || { echo '作業ツリーが dirty です。コミットしてから make release してください。' >&2; exit 1; }; \
	VERSION=$$($(MAKE) -s computed-version KIND="$(KIND)"); \
	echo "リリース: $$VERSION"; \
	if git rev-parse "refs/tags/$$VERSION" >/dev/null 2>&1; then echo "タグ $$VERSION はすでにあります。" >&2; exit 1; fi; \
	$(MAKE) test; \
	$(MAKE) build VERSION=$$VERSION GOOS=$(RELEASE_GOOS) GOARCH=$(RELEASE_GOARCH); \
	if [ -z "$$CI" ]; then git push origin HEAD; fi; \
	git tag "$$VERSION"; \
	git push origin "refs/tags/$$VERSION"; \
	gh release create "$$VERSION" --title "$$VERSION" --generate-notes "$(BINARY)"
