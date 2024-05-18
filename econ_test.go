package teeworldsecon

import (
	"testing"
)

const (
	econPort     = 7000
	econPassword = "hello_world"
)

func defaultEcon() *Econ {
	return NewEcon(
		NewDefaultEconConfig(
			econPort,
			econPassword,
		),
	)
}

func TestEconConnect(t *testing.T) {
	econ := defaultEcon()

	if err := econ.Connect(); err != nil {
		t.Error(err)
	}

	if err := econ.Auth(); err != nil {
		t.Error(err)
	}
}
