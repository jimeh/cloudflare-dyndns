DEV_DEPS = github.com/kardianos/govendor \
github.com/vektra/mockery/.../

BINARY = bin/cloudflare-dyndns
BINDIR = $(shell dirname ${BINARY})
SOURCES = $(shell find . -name '*.go')

.DEFAULT_GOAL: $(BINARY)
$(BINARY): $(SOURCES)
	go build -o ${BINARY}

.PHONY: build
build: $(BINARY)

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi; \
	if [ -d ${BINDIR} ]; then rmdir ${BINDIR}; fi

.PHONY: run
run: $(BINARY)
	$(BINARY)

.PHONY: install
install: dev-deps
	@govendor install +local +program

.PHONY: test
test: dev-deps
	@govendor test +local +program

.PHONY: vendor-sync
vendor-sync: dev-deps
	@govendor sync

.PHONY: vendor-fetch
vendor-fetch: dev-deps
	@govendor fetch +external +missing

.PHONY: vendor-install
vendor-install: dev-deps
	@govendor install +vendor

.PHONY: dev-deps
dev-deps:
	@$(foreach DEP,$(DEV_DEPS),go get $(DEP);)

.PHONY: update-dev-deps
update-dev-deps:
	@$(foreach DEP,$(DEV_DEPS),go get -u $(DEP);)

.PHONY: package
package:
	./package.sh
