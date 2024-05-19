package teeworldsecon

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"time"
)

const (
	EconPasswordMessage    = "Enter password:"
	EconAuthSuccessMessage = "Authentication successful. External console access granted."
	EconAuthFailMessage    = "Wrong password "
	EconResponseDuration   = 5
)

// Represents a Econ response
type EconResponse struct {
	// Raw value
	Value string
	// Indicates the response state, true if success
	State bool
}

// Econ client controller
type Econ struct {
	// Econ server configuration
	config *EconConfig
	// TCP Socket
	conn *net.Conn
	// Event channel
	eventCh chan string
	// Response channel
	responseCh chan EconResponse
	// Event manager
	EventManager *EconEventManager
}

// Create a Econ struct
func NewEcon(config *EconConfig) *Econ {
	return &Econ{
		config:       config,
		conn:         nil,
		eventCh:      make(chan string),
		responseCh:   make(chan EconResponse),
		EventManager: NewEconEventManager(),
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

	econ.conn = &conn

	return nil
}

// Goroutine for event listening
func (econ *Econ) listenEvents(errCh chan error) {
	if econ.conn == nil {
		errCh <- fmt.Errorf("missing connection")
		return
	}

	errCh <- nil

	scanner := bufio.NewScanner(*econ.conn)

	for scanner.Scan() {
		line := scanner.Text()
		econ.eventCh <- line
	}
}

// Start listening events
func (econ *Econ) ListenEvents() error {
	errCh := make(chan error)

	go econ.listenEvents(errCh)

	err := <-errCh

	if err != nil {
		return err
	}

	return nil
}

// The event manager calls the functions mapped with certain events
func (econ *Econ) HandleEvents() {
	for {
		data := <-econ.eventCh

		Debug(data)

		econ.EventManager.Call(data)
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
func (econ *Econ) WaitResponse(
	successMessage string,
	failMessage string,
) (*EconResponse, error) {
	if econ.conn == nil {
		return nil, fmt.Errorf("missing connection")
	}

	errCh := make(chan error)

	go func(errCh chan error) {
		var line string

		response := EconResponse{}
		scanner := bufio.NewScanner(*econ.conn)
		found := false

		for scanner.Scan() {
			line = scanner.Text()

			ok, _ := regexp.MatchString(failMessage, line)
			if ok {
				response.State = false
				found = true
				break
			}

			ok, _ = regexp.MatchString(successMessage, line)
			if ok {
				response.State = true
				found = true
				break
			}
		}

		response.Value = line

		if err := scanner.Err(); err != nil {
			errCh <- err
			return
		}

		if !found {
			errCh <- fmt.Errorf("cannot get an acceptable response")
			return
		}

		econ.responseCh <- response
	}(errCh)

	select {
	case response := <-econ.responseCh:
		return &response, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(EconResponseDuration * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// Send a payload then wait for its response
func (econ *Econ) SendAndWaitResponse(
	payload string,
	successMessage string,
	failMessage string,
) (*EconResponse, error) {
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

func (econ *Econ) Reconnect() error {
	Debug("waiting for %s", econ.address())

	for {
		err := econ.Connect()
		if err == nil {
			break
		}
	}

	return nil
}
