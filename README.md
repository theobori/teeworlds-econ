# Teeworlds econ server library

![golangci-lint](https://github.com/theobori/teeworlds-econ/actions/workflows/lint.yml/badge.svg)

This library is highly flexible and thread-safe by design, it allows you to interact with a Teeworlds econ server. It provides tools and help for abstracting econ server functionalities, such as in-game command error handling or event handling.

It is a kind of meta library.

## 📖 Build and run

You only need the following requirements:

- [Go](https://golang.org/doc/install) 1.22.3

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## 🧪 Tests

There are some tests that require a running Teeworlds econ server, feel free to use the `econ_server.sh` script that create a DDNet server and a econ server.

It also requires the following environment variables.

| Name | Description | Optional
| - | - | - |
`ECON_DEBUG` | Enables the debug verbose | Yes
`ECON_PORT` | Specifies the econ port | Yes (7000 by default)
`ECON_PASSWORD` | Specifies the econ password | No

By default, the test are using `localhost` as host.

Below, an example of running the test.

```bash
# Override some variables for the container and the Go tests
export ECON_PORT=1234
export ECON_PASSWORD=just_a_test_password

# Run the script
./econ_test.sh

# Wait for the econ server being ready
sleep 10

# Run the Go tests
make test
```

## 📎 Some examples

You can also take a look at [`command_kick.go`](./command_kick.go) and [`command_ban.go`](./command_ban.go) command

### Authenticate to the econ server

```go
package main

import twecon "github.com/theobori/teeworlds-econ"

func main() {
    // Econ server configuration
    config := twecon.EconConfig{
        Host: "127.0.0.1",
        Port: 7000,
        Password: "hello_world",
    }

    // Create the econ controller
    econ := twecon.NewEcon(&config)

    // Connect to the econ server
    if err := econ.Connect(); err != nil {
        return
    }

    // Authenticate to the econ server
    if _, err := econ.Authenticate(); err != nil {
        return
    }
}
```

### Send a payload

```go
package main

func main {
    ...
    // Once you are authenticated

    // Send a payload to the econ server,
	// in this case it is sending the broadcast command and its argument
	err := econ.Send("broadcast server shutdown in 5 minutes")
	if err != nil {
		return
	}
}
```

### Create your own command

```go
package main

import (
	"fmt"

	twecon "github.com/theobori/teeworlds-econ"
)

func main() {
    ...
    // Once you are authenticated

    // Create a your own command
    // This command wraps the in-game `say` command
    // So, as in-game, it takes one string argument
	sayCommand := twecon.EconCommand{
		Name: "say",
		ArgumentsAmount: 1,
		Func: func(econ *twecon.Econ, arguments ...any) (*twecon.EconResponse, error) {
			message, ok := arguments[0].(string)
			if !ok {
				return nil, fmt.Errorf("invalid message type")
			}

			payload := "say From the server: " + message

			if err := econ.Send(payload); err != nil {
				return nil, err
			}

			return &twecon.EconResponse{State: true}, nil
		},
	}

	// Register your command
	err = econ.CommandManager.Register(&sayCommand)
	if err != nil {
		return
	}

	// Execute your registered command
	_, err = econ.CommandManager.Exec(
		econ,
		"say",
		"server shutdown in 1 minute",
	)
	if err != nil {
		return
	}
}
```

### Event handler for player messages

```go
package main

import twecon "github.com/theobori/teeworlds-econ"

func main() {
    ...
    // Once you are authenticated

    // Create an event handlers
	playerMessage := twecon.EconEvent{
		Name: "player_message",
		Regex: "chat: .*:.*: .*",
		Func: func(econ *twecon.Econ, eventPayload string) any {
			fmt.Println(eventPayload)

			return nil
		},
	}

	// Register the simple event
	err = econ.EventManager.Register(&playerMessage)
	if err != nil {
		return
	}
}
```

### A more complex event handler

```go
package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	twecon "github.com/theobori/teeworlds-econ"
)

func main() {
    ...
    // Once you are authenticated

    // Custom in-game chatbot
	chatbot := twecon.EconEvent{
		Name: "player_chat_bot",
		Regex: "chat: .*:.*: .*",
		Func: func(econ *twecon.Econ, eventPayload string) any {
			re := regexp.MustCompile(`chat: .*:(.*): (.*)`)
			matches := re.FindStringSubmatch(eventPayload)

			msg := matches[2]

			if strings.HasPrefix(msg, "!random") {
				_, err := econ.CommandManager.Exec(
					econ,
					"say",
					fmt.Sprintf("%d", rand.Intn(100)),
				)

				if err != nil {
					return nil
				}
			}

			return nil
		},
	}

	// Register the chatbot event
	err = econ.EventManager.Register(&chatbot)
	if err != nil {
		return
	}

    // Start handling the events
	// `HandleEvents` is a control loop, by default it is not a goroutine
	go econ.HandleEvents()
}
```
