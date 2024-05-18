package teeworldsecon

import "fmt"

// Teeworlds broadcast command
func (econ *Econ) Broadcast(payload string) error {
	return econ.Send(fmt.Sprintf("broadcast %s", payload))
}
