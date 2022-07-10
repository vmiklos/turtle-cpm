GO_OBJECTS = \
	cpm.go \
	cpm_test.go \
	main.go \

build: cpm

cpm: Makefile ${GO_OBJECTS}
	go build

check: build check-format check-lint
	@echo "make check: ok"

check-lint:
	golint -set_exit_status

check-format:
	[ -z "$(shell gofmt -l ${GO_OBJECTS})" ]
