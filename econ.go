package teeworldsecon

import (
	"fmt"
	"net"
	"strings"
)

const (
	EconPasswordMessage        = "Enter password:"
	EconAuthSuccessfullMessage = "Authentication successful. External console access granted."
	EconAuthWrongMessage       = "Wrong password "
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

// Authenticate to the econ server
func (econ *Econ) Auth() error {
	var bufString string

	if econ.conn == nil {
		return fmt.Errorf("missing connection")
	}

	buf := make([]byte, 256)
	conn := *econ.conn

	_, err := conn.Write([]byte(econ.config.Password + "\n"))
	if err != nil {
		return err
	}

	for {
		_, err := conn.Read(buf)
		if err != nil {
			return err
		}

		bufString = string(buf)

		if strings.HasPrefix(bufString, EconAuthSuccessfullMessage) {
			break
		}

		if strings.HasPrefix(bufString, EconAuthWrongMessage) {
			return fmt.Errorf("authentication failed: %s", bufString)
		}
	}

	return nil
}

// Disconnect from the econ server
func (econ *Econ) Disconnect() error {
	if econ.conn == nil {
		return fmt.Errorf("missing connection")
	}

	return (*econ.conn).Close()
}
