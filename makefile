GREEN=\n\033[1;32;40m
NC=\033[0m # No Color

mocks:
	@/bin/sh -c 'echo "${GREEN}[테스트를 시작합니다.]${NC}"'
	@mockery --dir modules --all --case underscore --keeptree
.PHONY: mocks