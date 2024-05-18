package teeworldsecon

import "fmt"

// Teeworlds say command
func (econ *Econ) Say(payload string) error {
	return econ.Send(fmt.Sprintf("say %s", payload))
}
