package logger

import (
	"log"
	"os"
)

func Init() {
	log.SetOutput(os.Stdout)
}
