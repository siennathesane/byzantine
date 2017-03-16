## Byzantine

[![license](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square)]()
[![GoDoc](https://godoc.org/github.com/mxplusb/byzantine?status.svg)](https://godoc.org/github.com/mxplusb/byzantine)
[![Go Report Card](https://goreportcard.com/badge/github.com/mxplusb/byzantine)](https://goreportcard.com/report/github.com/mxplusb/byzantine)

This is an implementation of the Byzantine Fault Tolerance algorithm as modelled by [Eric Scott Freeman](https://brage.bibsys.no/xmlui/bitstream/handle/11250/2413908/Freeman_Eric.pdf?sequence=1&isAllowed=y) with gRPC for Go.

### Generating Changes

To generate changes with `protoc`:

```protoc -I/usr/local/include -I. \
     -I$GOPATH/src \
     -I$GOPATH/src/github.com/googleapis/googleapis/ \
     -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
     --go_out=,plugins=grpc:. \
     --swagger_out=logtostderr=true:. \
     --grpc-gateway_out=logtostderr=true:. \
     byzantine.proto```