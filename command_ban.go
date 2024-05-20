package teeworldsecon

import "fmt"

var (
	EconBanSuccessMessage = `net_ban: banned '.*' for \d+ minute`
	EconBanFailMessage    = `net_ban: ban error`
)

// DDNet ban command
var DDNetBan = func(econ *Econ, arguments ...any) (*EconResponse, error) {
	// IP or ID
	player, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid player type")
	}

	minutes, ok := arguments[1].(int)
	if !ok {
		return nil, fmt.Errorf("invalid minutes type")
	}

	reason, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reason type")
	}

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
