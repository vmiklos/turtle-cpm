GO_OBJECTS = \
	commands/create.go \
	commands/create_test.go \
	commands/delete.go \
	commands/delete_test.go \
	commands/import.go \
	commands/import_test.go \
	commands/read.go \
	commands/read_test.go \
	commands/root.go \
	commands/root_test.go \
	commands/sync.go \
	commands/sync_test.go \
	commands/update.go \
	commands/update_test.go \
	commands/version.go \
	commands/version_test.go \
	main.go \
	man/generate.go \

COMMANDS_PATH = vmiklos.hu/go/cpm/commands
LDFLAGS = -X ${COMMANDS_PATH}.GitCommit=$(shell git describe)

build: turtle-cpm

turtle-cpm: Makefile ${GO_OBJECTS}
	go build -ldflags "$(LDFLAGS)" .

check: build check-format check-lint check-unit
	@echo "make check: ok"

check-lint:
	golint -set_exit_status ./...

check-format:
	[ -z "$(shell gofmt -l ${GO_OBJECTS})" ]

# Without coverage: 'go test ./...'.
check-unit:
	go build -ldflags "$(LDFLAGS)" ./...
	courtney -e ./...

generate-man:
	go run man/generate.go

run-guide:
	cd guide && mdbook serve --hostname 127.0.0.1

codespell:
	codespell $(shell git ls-files)
