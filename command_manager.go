package teeworldsecon

import (
	"fmt"
	"sync"
)

// Command manager controller
type EconCommandManager struct {
	commands map[string]*EconCommand
	mu       sync.Mutex
}

// Create a command manager
func NewEconCommandManager() *EconCommandManager {
	return &EconCommandManager{
		commands: make(map[string]*EconCommand),
	}
}

// Returns a command entry for the commands registry
func (em *EconCommandManager) Entry(name string) *EconCommand {
	em.mu.Lock()
	defer em.mu.Unlock()

	return em.commands[name]
}

// Delete a command
func (em *EconCommandManager) Delete(name string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	delete(em.commands, name)
}

// Register a command
func (em *EconCommandManager) Register(command *EconCommand) error {
	if command == nil {
		return fmt.Errorf("nil command")
	}

	em.mu.Lock()
	defer em.mu.Unlock()

	em.commands[command.Name] = command

	return nil
}

// Execute a command if exists
func (em *EconCommandManager) Exec(econ *Econ, name string, arguments ...any) (*EconResponse, error) {
	if econ == nil {
		return nil, fmt.Errorf("econ is nil")
	}

	em.mu.Lock()
	defer em.mu.Unlock()

	command, found := em.commands[name]
	if !found {
		return nil, fmt.Errorf("command does not exist")
	}

	if len(arguments) != command.ArgumentsAmount {
		return nil, fmt.Errorf(
			"%s requires %d",
			command.Name,
			command.ArgumentsAmount,
		)
	}

	return command.Func(econ, arguments...)
}

// Returns every command names
func (em *EconCommandManager) Names() []string {
	ret := []string{}

	for name := range em.commands {
		ret = append(ret, name)
	}

	return ret
}
