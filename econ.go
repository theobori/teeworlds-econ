package teeworldsecon

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/theobori/teeworlds-econ/internal"
)

const (
	EconPasswordMessage    = "Enter password:"
	EconAuthSuccessMessage = "Authentication successful. External console access granted."
	EconAuthFailMessage    = "Wrong password "
	EconResponseDuration   = 5
)

// Econ client controller
type Econ struct {
	// Server
	// Econ server configuration
	config *EconConfig
	// TCP Socket
	conn *net.Conn

	// Managers
	// Event manager
	EventManager *EconEventManager
	// Command manager
	CommandManager *EconCommandManager
	// Reponse manager
	reponseManager *EconResponseManager
	// Payload manager
	payloadManager *EconResponseManager
}

// Create a Econ struct
func NewEcon(config *EconConfig) *Econ {
	return &Econ{
		config:         config,
		conn:           nil,
		EventManager:   NewEconEventManager(),
		CommandManager: NewEconCommandManager(),
		reponseManager: NewEconResponseManager(),
		payloadManager: NewEconResponseManager(),
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

	// Set deadline
	deadline := time.Now().Add(EconResponseDuration * time.Second)
	err = conn.SetDeadline(deadline)
	if err != nil {
		return err
	}

	buf := make([]byte, len(EconPasswordMessage))

	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	// Remove deadline
	err = conn.SetDeadline(time.Time{})
	if err != nil {
		return err
	}

	if string(buf) != EconPasswordMessage {
		return fmt.Errorf("invalid econ server")
	}

	// Set the connection
	econ.conn = &conn

	// Start listening events
	go econ.listenEvents()

	return nil
}

// Listening the server events
func (econ *Econ) listenEvents()  {
	if econ.conn == nil {
		teeworldsecon.Debug("missing connection")
	}

	scanner := bufio.NewScanner(*econ.conn)

	for scanner.Scan() {
		line := scanner.Text()
		// Send to the event channels if needed
		econ.payloadManager.Send(line)
		// Send to the reponse channels if needed
		econ.reponseManager.Send(line)
	}

	err := scanner.Err()
	if err != nil {
		teeworldsecon.Debug("%v", err)
	}
}

// Control loop for the events, calling function
// registered with the event manager
// that match the current event
func (econ *Econ) HandleEvents() {
	eventCh := make(chan string, 1)

	econ.payloadManager.Add(eventCh)

	for {
		data := <-eventCh

		econ.EventManager.Handle(econ, data)
	}
}

// Send a payload to the Econ server
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

// Wait for a server response
func (econ *Econ) WaitResponse(successMessage string, failMessage string) (*EconResponse, error) {
	if econ.conn == nil {
		return nil, fmt.Errorf("missing connection")
	}

	responseCh := make(chan EconResponse, 1)
	payloadCh := make(chan string, 1)

	id := econ.reponseManager.Add(payloadCh)

	go func(
		payloadCh chan string,
		responseCh chan EconResponse,
	) {
		var line string

		response := EconResponse{}

		for {
			line = <-payloadCh

			ok, _ := regexp.MatchString(failMessage, line)
			if ok {
				response.State = false
				break
			}

			ok, _ = regexp.MatchString(successMessage, line)
			if ok {
				response.State = true
				break
			}
		}

		response.Value = line

		responseCh <- response
	}(payloadCh, responseCh)

	defer econ.reponseManager.Delete(id)

	select {
	case response := <-responseCh:
		return &response, nil
	case <-time.After(EconResponseDuration * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// Send a payload then wait for its response
func (econ *Econ) SendAndWaitResponse(payload string, successMessage string, failMessage string) (*EconResponse, error) {
	err := econ.Send(payload)
	if err != nil {
		return nil, err
	}

	response, err := econ.WaitResponse(
		successMessage,
		failMessage,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Authenticate to the econ server
func (econ *Econ) Authenticate() (*EconResponse, error) {
	return econ.SendAndWaitResponse(
		econ.config.Password,
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
