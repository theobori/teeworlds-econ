package teeworldsecon

import (
	"strings"
	"testing"
)

// Testing the econ event manager
func TestEconEventManager(t *testing.T) {
	em := NewEconEventManager()

	playerChat := NewEconEvent(
		"playerChat",
		"chat: .*",
		func(eventPayload string) any {
			return strings.Split(eventPayload, ": ")[1]
		},
	)

	em.Register(playerChat)

	em.Call("chat: hello_world1")
	em.Call("chat: hello_world2")

	if len(playerChat.Queue) != 2 {
		t.Errorf("Wrong number of event handled")
	}
}
