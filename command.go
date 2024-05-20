package teeworldsecon

// Function prototype for a command
type EconCommandFunc func(econ *Econ, arguments ...any) (*EconResponse, error)

// Represents a command
type EconCommand struct {
	// Command name
	Name string
	// Command arguments amount
	ArgumentsAmount int
	// Command function
	Func EconCommandFunc
}
