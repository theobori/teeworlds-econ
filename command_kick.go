package teeworldsecon

import "fmt"

var (
	EconKickFailMessage = "server: invalid client id to kick"
)

// Teeworlds kick command
func (econ *Econ) Kick(id uint8, reason string) (*EconResponse, error) {
	var m string

	payload := fmt.Sprintf("kick %d", id)

	if reason != "" {
		payload += " " + reason
		m = fmt.Sprintf(`Kicked \(%s\)`, reason)
	} else {
		m = "Kicked by console"
	}

	successMessage := fmt.Sprintf(`\(%s\)`, m)

	return econ.SendAndWaitResponse(
		payload,
		successMessage,
		EconKickFailMessage,
	)
}
