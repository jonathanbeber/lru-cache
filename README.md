# lru-cache

**L**east **R**ecently **U**sed cache implementation using Go. It's a proof of concept and shouldn't be used in production environments.

## Using it

It exports one struct `github.com/jonathanbeber/lru-cache/lru`.`Cache`. Its constructor receives the maximum number of items to cache, a `WrappableFunction` and a `Logger`. Check the code documentation for more details (`godoc`). After creating a `lru.Cache` object, eg `lruc`, invoke the wrapped function with `lruc.Do(key)`.

For better understanding, read the example code in `./example` and run the factorial one with the command `make factorial`.

# Development

## Pre-requisites

This application is developed to go 1.14. Check the [official docs](https://golang.org/doc/install) for details on how to install it.
It uses [golangci-lint][0] and [golint][1] for code analyzes. It'll install them when running `make lint` or `make test` if the commands are not found.

## Commands

The Makefile contains some commands for better productivity:

* To test the application use `make test`. It will run all the unit tests available. It will also run the following commands, available as separated `make` commands:
  * `make fmt`: format all files using `gofmt` tool. It'll apply changes to the files. This project aggress on ceding control over minutiae of hand-formatting in favour of the `gofmt` tool result;
  * `make vet`: examines Go source code and reports suspicious constructs;
  * `make lint`:  runs the [golangci-lint][0] and the [golint][1] tool;
* `make factorial` run the code in `./examples/factorial`.

[0]: https://github.com/golangci/golangci-lint
[1]: https://github.com/golang/lint
