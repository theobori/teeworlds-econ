package teeworldsecon

import "fmt"

var (
	EconBanSuccessMessage = `net_ban: banned '.*' for \d+ minute`
	EconBanFailMessage    = `net_ban: ban error`
)

// Teeworlds ban command
func (econ *Econ) Ban(player string, minutes int, reason string) (*EconResponse, error) {
	payload := fmt.Sprintf(
		"ban %s %d %s",
		player,
		minutes,
		reason,
	)

	return econ.SendAndWaitResponse(
		payload,
		EconBanSuccessMessage,
		EconBanFailMessage,
	)
}
