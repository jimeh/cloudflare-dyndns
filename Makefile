DEV_DEPS = github.com/kardianos/govendor \
github.com/vektra/mockery/.../ \

BIN_PATH = ./bin/cloudflare-dyndns

test: dev-deps
	@govendor test +local +program

install: dev-deps
	@govendor install +local +program

build:
	mkdir -p bin && go build -o $(BIN_PATH)

run: build
	$(BIN_PATH)

fetch-vendor: dev-deps
	@govendor fetch +external +missing

install-vendor: dev-deps
	@govendor install +vendor

dev-deps:
	@$(foreach DEP,$(DEV_DEPS),go get $(DEP);)

update-dev-deps:
	@$(foreach DEP,$(DEV_DEPS),go get -u $(DEP);)

.PHONY: test install build run fetch-vendor install-vendor dev-deps \
	update-dev-deps
