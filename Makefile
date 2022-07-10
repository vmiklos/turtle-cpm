GO_OBJECTS = \
	main.go \
	cpm.go \

build: cpm

cpm: Makefile ${GO_OBJECTS}
	go build

check: build check-format check-lint
	@echo "make check: ok"

check-lint:
	golint -set_exit_status

check-format:
	[ -z "$(shell gofmt -l ${GO_OBJECTS})" ]
