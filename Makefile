GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
RELEASE_OPTIONS := -ldflags "-s -w -X github.com/goby-lang/goby/vm.DefaultLibPath=${GOBY_LIBPATH}" -tags release
TEST_OPTIONS := -ldflags "-s -w"

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: build
build:
	go build $(RELEASE_OPTIONS) .

# If $GOBY_BINPATH is not specified, the standard go install is performed.
# Otherwise, the binary is installed there, and the `lib` directory copied as `$GOBY_LIBPATH`.
#
.PHONY: install
install:
ifeq ($(GOBY_BINPATH),)
	go install $(RELEASE_OPTIONS) .
else
ifndef GOBY_LIBPATH
	$(error GOBY_BINPATH requires GOBY_LIBPATH to be set)
endif
	go build $(RELEASE_OPTIONS) ./src/github.com/goby-lang/goby
	mkdir -p ${DESTDIR}${GOBY_BINPATH}
	mv -v goby ${DESTDIR}${GOBY_BINPATH}
	mkdir -p ${DESTDIR}${GOBY_LIBPATH}
	cp -R src/github.com/goby-lang/goby/lib ${DESTDIR}${GOBY_LIBPATH}
endif

.PHONY: test
test:
	go test $(TEST_OPTIONS) ./...

.PHONY: clean
clean:
	go clean .
