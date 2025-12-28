package utils

import (
	"log"
	"os"
)

func InitLogger() {
	// Initialize logger (simple wrapper for now)
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
