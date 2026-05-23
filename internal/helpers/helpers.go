package helpers

import (
	"log"
)

func Die(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func Todo(msg string) {
	panic(msg)
}
