package teeworldsecon

import "fmt"

var (
	EconKickFailMessage = "server: invalid client id to kick"
)

// DDNet kick command
var DDNetKick = func(econ *Econ, arguments ...any) (*EconResponse, error) {
	var m string

	id, ok := arguments[0].(int)
	if !ok {
		return nil, fmt.Errorf("invalid id type")
	}

	reason, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid reason type")
	}

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
