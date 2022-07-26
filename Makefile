GO_OBJECTS = \
	commands/commands_test.go \
	commands/create.go \
	commands/delete.go \
	commands/import.go \
	commands/read.go \
	commands/root.go \
	commands/update.go \
	main.go \
	man/generate.go \

build: turtle-cpm

turtle-cpm: Makefile ${GO_OBJECTS}
	go build .

check: build check-format check-lint check-unit
	@echo "make check: ok"

check-lint:
	golint -set_exit_status ./...

check-format:
	[ -z "$(shell gofmt -l ${GO_OBJECTS})" ]

# Without coverage: 'go test ./...'.
check-unit:
	go build ./...
	courtney -e ./...

generate-man:
	go run man/generate.go

run-guide:
	cd guide && mdbook serve --hostname 127.0.0.1
