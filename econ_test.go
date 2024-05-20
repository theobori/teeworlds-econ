package teeworldsecon

import (
	"os"
	"strconv"
	"testing"
)

const (
	econPort = 7000
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
		return econ
	}

	if _, err := econ.Authenticate(); err != nil {
		t.Error(err)
		return econ
	}

	return econ
}

// Testing the econ server connection
func TestEconConnect(t *testing.T) {
	econ := defaultEcon()

	econConnectAndAuth(econ, t)
}

// Testing that we receive an error
func TestEconDDNetKick(t *testing.T) {
	econ := econConnectAndAuth(defaultEcon(), t)

	command := EconCommand{
		Name:            "ddnet_kick",
		ArgumentsAmount: 2,
		Func:            DDNetKick,
	}

	if err := econ.CommandManager.Register(&command); err != nil {
		t.Error(err)
	}

	r, err := econ.CommandManager.Exec(econ, "ddnet_kick", 0, "abuse")
	if err != nil || r.State {
		t.Error(err)
	}
}

// Testing that we receive an error
func TestEconDDNetBan(t *testing.T) {
	econ := econConnectAndAuth(defaultEcon(), t)

	command := EconCommand{
		Name:            "ddnet_ban",
		ArgumentsAmount: 3,
		Func:            DDNetBan,
	}

	if err := econ.CommandManager.Register(&command); err != nil {
		t.Error(err)
	}

	r, err := econ.CommandManager.Exec(econ, "ddnet_ban", "0", 10, "block")
	if err != nil || r.State {
		t.Error(err)
	}
}
