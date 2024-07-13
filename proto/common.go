//go:generate go run github.com/webrpc/webrpc/cmd/webrpc-gen -schema=proto.ridl -target=golang@v0.14.8 -pkg=proto -server -client -out=./proto.gen.go
package proto

import (
	_ "github.com/webrpc/webrpc"
)
