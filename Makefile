GO		?= go
OS      := $(shell uname)
ARCH	:= $(shell uname -m)
VERSION	:= $(shell git rev-parse HEAD)

TAGS	:= netgo
GOFLAGS	:= -installsuffix netgo
LDFLAGS	:= -extldflags "-static"

TESTS	:= .*

TARGET	:= kexp

all: build

.PHONY: build
build:
build: LDFLAGS += -X main.version='$(VERSION)'
build:
	@echo "BUILD $(TARGET)"
	@$(GO) build -tags '$(TAGS)' $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(TARGET) $^

.PHONY: clean
clean:
	@echo "[CLEAN] $(TARGET)"
	@rm -rf $(TARGET)


.PHONY: test
test:
	@$(GO) test -cover -tags '$(TAGS)' -run '$(TESTS)'

