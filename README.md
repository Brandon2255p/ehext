[![Build Status](https://travis-ci.org/Brandon2255p/ehext.svg?branch=master)](https://travis-ci.com/Brandon2255p/ehext)
[![Coverage Status](https://img.shields.io/coveralls/Brandon2255p/ehext.svg)](https://coveralls.io/r/Brandon2255p/ehext)
[![GoDoc](https://godoc.org/github.com/Brandon2255p/ehext?status.svg)](https://godoc.org/github.com/Brandon2255p/ehext)
[![Go Report Card](https://goreportcard.com/badge/Brandon2255p/ehext)](https://goreportcard.com/report/Brandon2255p/ehext)

# ehext
Extensions for the go package github.com/looplab/eventhorizon

## Scaffolding Generator
The ehext generator will read command and event structs and generate the domain scaffolding needed to use eventhorizon

To install the generator:
```
go get github.com/Brandon2255p/ehext
```
```
go install  github.com/Brandon2255p/ehext
```

## Example
- Install the generator 
- Download the example code
- `cd ehext/example/domain`
- `go generate`