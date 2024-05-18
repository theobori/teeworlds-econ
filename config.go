package teeworldsecon

// Teeworlds econ server configuration
type EconConfig struct {
	// Server IP address
	Host string
	// Server port
	Port uint16
	// Server password
	Password string
}

func NewDefaultEconConfig(port uint16, password string) *EconConfig {
	return &EconConfig{
		Host:     "localhost",
		Port:     port,
		Password: password,
	}
}
