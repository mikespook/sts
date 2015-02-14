package rpc

import (
	"testing"

	"github.com/mikespook/sts/bus"
)

func TestInterface(t *testing.T) {
	var service bus.Service = &RPC{}
	t.Logf("%T", service)
}
