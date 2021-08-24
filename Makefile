BIN = moac
CGO_ENABLED ?= 0
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin

GO ?= go
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOKART ?= $(GOBIN)/gokart
GOKART_FLAGS ?= -g

# general build flags
GO_BUILDFLAGS += -trimpath -mod=readonly -modcacherw

# if building with CGO, turn on some hardening
CC ?= clang
CCLD ?= lld
CFLAGS ?= -O2 -fno-semantic-interposition -g -pipe -Wp,-D_FORTIFY_SOURCE=2 -Wp,-D_GLIBCXX_ASSERTIONS -fexceptions -fstack-protector-all -m64 -fasynchronous-unwind-tables -fstack-clash-protection -fcf-protection=full -ffunction-sections -fdata-sections -ftrivial-auto-var-init=zero -enable-trivial-auto-var-init-zero-knowing-it-will-be-removed-from-clang
LDFLAGS += -Wl,-z,relro -Wl,--as-needed -Wl,-z,now -Wl,-E -Wl,-z,noexecstack -Wl,--gc-sections
GO_LDFLAGS += "-w -s -linkmode=external -extldflags '$(LDFLAGS)'"

# for testing with clang+msan+CFI+safe-stack/shadow-stack and release builds
CFLAGS_CFI += $(CFLAGS) -flto=thin -fvisibility=hidden -fsanitize=cfi -fpic -fpie
LDFLAGS_CFI += $(LDFLAGS) -flto=thin -fsanitize=cfi -pie
GO_LDFLAGS_CFI += "-w -s -linkmode=external -extldflags '$(LDFLAGS_CFI)'"

# for release builds, with Clang+CFI sanitization, static-pie linked
LDFLAGS_RELEASE += $(LDFLAGS) -flto=thin -fsanitize=cfi,safe-stack -static-pie
GO_LDFLAGS_RELEASE += "-w -s -linkmode=external -extldflags '$(LDFLAGS_RELEASE)'"

default:
	$(MAKE) build

lint:
	@echo "LINTING"
	$(GOLANGCI_LINT) run
	$(GOKART) scan $(GOKART_FLAGS) .
	$(GOKART) scan $(GOKART_FLAGS) ./entropy
	$(GOKART) scan $(GOKART_FLAGS) ./cmd/moac

# Test with thread and memory sanitizers; needs associated libclang_rt libs.
# `make test` does not work on alpine (compiler-rt lacks msan)
# but it works on fedora and void-musl.
test:
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS)" $(GO) test -race -ldflags=$(GO_LDFLAGS)
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) test $(GO_BUILDFLAGS) -buildmode=pie -msan -ldflags=$(GO_LDFLAGS_CFI)

$(BIN):
	CC=$(CC) CCLD=$(CCLD) CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GO_BUILDFLAGS) -o $(BIN) ./cmd/moac/

build: $(BIN)

clean:
	$(GO) clean

build-safe:
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) build $(GO_BUILDFLAGS) -buildmode=pie -ldflags=$(GO_LDFLAGS_CFI) -o $(BIN) ./cmd/moac

# build-release builds a static-pie binary with sanitizers for CFI and either
# safe-stack (x86_64) or shadow-call-stack (ARMv8)
# the below should be run on a musl-based toolchain; works on Alpine or Void-musl
# Tends to cause crashes when linking with glibc
build-release:
	CC=$(CC) CCLD=$(CCLD) CGO_CFLAGS="$(CFLAGS_CFI)" $(GO) build $(GO_BUILDFLAGS) -buildmode=pie -ldflags=$(GO_LDFLAGS_RELEASE) -o $(BIN) ./cmd/moac

.PHONY: all lint test build build-release build-safe clean
