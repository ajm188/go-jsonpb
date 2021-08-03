# go-jsonpb

[![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GoDoc](https://godoc.org/github.com/ajm188/go-jsonpb?status.svg)](http://godoc.org/github.com/ajm188/go-jsonpb)
[![Go Report Card](https://goreportcard.com/badge/github.com/ajm188/go-jsonpb)](https://goreportcard.com/report/github.com/ajm188/go-jsonpb)

Add protobuf/jsonpb marshalers to your proto messages

## Usage

To use:

1. Install the plugin: `go install github.com/ajm188/go-jsonpb/cmd/protoc-gen-go-json`.
    * Ensure that `protoc-gen-go-json` is in your `$PATH`. This should be the
case as long as `$GOPATH/bin` is in your `$PATH`.
2. Run protoc with the plugin: `protoc --go_out=. --go_json_out=. /path/to/my.proto`.
