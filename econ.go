package teeworldsecon

import (
	"fmt"
	"net"
	"regexp"
	"time"
)

const (
	// Server TCP messages

	EconPasswordMessage    = "Enter password:"
	EconAuthSuccessMessage = "Authentication successful. External console access granted."
	EconAuthFailMessage    = "Wrong password "
	EconKickFailMessage    = "server: invalid client id to kick"

	EconBanSuccessMessage = `net_ban: banned '.*' for \d+ minute`
	EconBanFailMessage = `net_ban: ban error`

	// Server offsets before messages

	EconBaseOffset   = 22
	EconServerOffset = EconBaseOffset + 8

	// Server TCP connection specifications
	EconServerDuration = 5
)

// Econ client controller
type Econ struct {
	// Econ server configuration
	config *EconConfig
	// TCP Socket
	conn *net.Conn
}

// Create a Econ struct
func NewEcon(config *EconConfig) *Econ {
	return &Econ{
		config: config,
		conn:   nil,
	}
}

// Set the config value
func (econ *Econ) SetEconConfig(config *EconConfig) {
	econ.config = config
}

// Return a formatted address of format 'host:port'
func (econ *Econ) address() string {
	return fmt.Sprintf(
		"%s:%d",
		econ.config.Host,
		econ.config.Port,
	)
}

// Connect to the econ server and check its validity
func (econ *Econ) Connect() error {
	address := econ.address()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	buf := make([]byte, len(EconPasswordMessage))

	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	if string(buf) != EconPasswordMessage {
		return fmt.Errorf("invalid econ server")
	}

	econ.conn = &conn

	return nil
}

func (econ *Econ) Send(payload string) error {
	if econ.conn == nil {
		return fmt.Errorf("missing connection")
	}

	_, err := (*econ.conn).Write([]byte(payload + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (econ *Econ) readWithTimeout(b []byte) error {
	if econ.conn == nil {
		return fmt.Errorf("missing connection")
	}

	deadline := time.Now().Add(EconServerDuration * time.Second)

    err := (*econ.conn).SetReadDeadline(deadline);
	if err != nil {
        return fmt.Errorf("failed to set read deadline: %v", err)
    }

	_, err = (*econ.conn).Read(b)
	if err != nil {
		return err
	}
	
	return nil
}

func (econ *Econ) waitResponse(
	successMessage string,
	failMessage string,
) error {
	var s string

	b := make([]byte, 256)

	for {
		err := econ.readWithTimeout(b)
		if err != nil {
			return err
		}

		s = string(b)

		ok, _ := regexp.MatchString(failMessage, s)
		if ok {
			return fmt.Errorf(s)
		}

		ok, _ = regexp.MatchString(successMessage, s)
		if ok {
			break
		}
	}

	return nil
}

// Authenticate to the econ server
func (econ *Econ) Auth() error {
	if econ.conn == nil {
		return fmt.Errorf("missing connection")
	}

	err := econ.Send(econ.config.Password)
	if err != nil {
		return err
	}

	return econ.waitResponse(
		EconAuthSuccessMessage,
		EconAuthFailMessage,
	)
}

// Disconnect from the econ server
func (econ *Econ) Disconnect() error {
	if econ.conn == nil {
		return fmt.Errorf("missing connection")
	}

	return (*econ.conn).Close()
}

// Teeworlds say command
func (econ *Econ) Say(payload string) error {
	return econ.Send(fmt.Sprintf("say %s", payload))
}

// Teeworlds broadcast command
func (econ *Econ) Broadcast(payload string) error {
	return econ.Send(fmt.Sprintf("broadcast %s", payload))
}

// Teeworlds kick command
func (econ *Econ) Kick(id uint8, reason string) error {
	var m string

	payload := fmt.Sprintf("kick %d", id)

	if reason != "" {
		payload += " " + reason
		m = fmt.Sprintf(`Kicked \(%s\)`, reason)
	} else {
		m = "Kicked by console"
	}

	m = fmt.Sprintf(`\(%s\)`, m)

	err := econ.Send(payload)
	if err != nil {
		return err
	}

	return econ.waitResponse(
		m,
		EconKickFailMessage,
	)
}

// Teeworlds ban command
func (econ *Econ) Ban(player string, minutes int, reason string) error {
	payload := fmt.Sprintf(
		"ban %s %d %s",
		player,
		minutes,
		reason,
	)

	err := econ.Send(payload)
	if err != nil {
		return err
	}

	return econ.waitResponse(
		EconBanSuccessMessage,
		EconBanFailMessage,
	)
}
