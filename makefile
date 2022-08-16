GREEN=\n\033[1;32;40m
NC=\033[0m # No Color

PKG_LIST := $(shell go list ./... | grep -v .back | grep -v config | grep -v pb)

mocks:
	@/bin/sh -c 'echo "${GREEN}[테스트를 시작합니다.]${NC}"'
	@mockery --dir modules --all --case underscore --keeptree
.PHONY: mocks

coverage:
	@/bin/sh -c 'echo "${GREEN}[test coverage를 계산합니다.]${NC}"'
	@mkdir -p .public/coverage
	@gocov test ${PKG_LIST} | gocov-html > .public/coverage/index.html
	@gocov test ${PKG_LIST} | gocov-xml > coverage.xml
	@gocov test ${PKG_LIST} | gocov report
.PHONY: coverage