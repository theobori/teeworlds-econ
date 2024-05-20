package teeworldsecon

import (
	"log"
	"os"
)

func Debug(format string, v ...any) {
	debug := os.Getenv("ECON_DEBUG")

	if debug == "1" {
		log.Printf(format, v...)
	}
}
