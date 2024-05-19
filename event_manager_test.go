package teeworldsecon

import (
	"strings"
	"testing"
)

// Testing the econ event manager
func TestEconEventManager(t *testing.T) {
	em := NewEconEventManager()

	playerChat := EconEvent{
		"playerChat",
		"chat: .*",
		func(eventPayload string) any {
			return strings.Split(eventPayload, ": ")[1]
		},
	}

	if err := em.Register(&playerChat); err != nil {
		t.Error(err)
	}

	if len(em.Names()) != 1 {
		t.Errorf("Invalid number of registered events")
	}

	em.Call("chat: hello_world1")
	em.Call("chat: hello_world2")

	entry := em.Entry("playerChat")
	if entry == nil {
		t.Error()
	}

	if entry.QueueSize() != 2 {
		t.Errorf("Wrong number of event handled")
	}
}
