# Development notes

## Updating dependencies

Ideally CI checks everything before a commit hits main, but run `go get -u && go mod tidy` from time
to time and make sure Go dependencies are reasonably up to date.

## Shell completion

Changes to the shell completion can be tested, without restarting the shell using:

```console
$ source <(cpm completion bash)
```
