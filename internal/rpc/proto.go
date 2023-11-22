package rpc

import "errors"

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto.proto

var ErrAgentDeactivated = errors.New("agent deactivated")
