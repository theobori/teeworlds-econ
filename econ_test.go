package teeworldsecon

import (
	"os"
	"strconv"
	"testing"
)

const (
	econPort     = 7000
	econPassword = "hello_world"
)

func defaultEcon() *Econ {
	value, err := strconv.ParseUint(os.Getenv("ECON_PORT"), 10, 16)

	if err != nil {
		value = econPort
	}

	return NewEcon(
		NewDefaultEconConfig(
			uint16(value),
			os.Getenv("ECON_PASSWORD"),
		),
	)
}

func econConnectAndAuth(econ *Econ, t *testing.T) *Econ {
	if err := econ.Connect(); err != nil {
		t.Error(err)
	}

	if _, err := econ.Auth(); err != nil {
		t.Error(err)
	}

	return econ
}

// Testing the econ server connection
func TestEconConnect(t *testing.T) {
	econ := defaultEcon()

	econConnectAndAuth(econ, t)
}

// Testing that we receive an error
func TestEconKick(t *testing.T) {
	econ := econConnectAndAuth(defaultEcon(), t)

	if r, err := econ.Kick(0, "reason"); err != nil || r.State {
		t.Error()
	}
}

// Testing that we receive an error
func TestEconBan(t *testing.T) {
	econ := econConnectAndAuth(defaultEcon(), t)

	if r, err := econ.Ban("3", 10, "reason"); err != nil || r.State {
		t.Error()
	}
}
