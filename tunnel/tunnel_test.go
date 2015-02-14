package tunnel

import (
	"testing"

	"github.com/mikespook/sts/bus"
)

func TestInterface(t *testing.T) {
	var service bus.Service = &Tunnel{}
	t.Logf("%T", service)
}

type _tc struct {
	c interface{}
	e bool
}

var _tcConfig = []_tc{
	_tc{c: struct{}{}, e: true},
	_tc{c: &Config{}, e: false},
}

func TestConfig(t *testing.T) {
	tunnel := &Tunnel{}
	for _, v := range _tcConfig {
		if err := tunnel.Config(v.c); (err == nil) == v.e {
			t.Errorf("Passing wrong paramater %T", v.c)
		}
	}
}
