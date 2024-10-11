package history

import (
	"log"
)

const timeoutSecond = 3

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s\n", err.Error(), msg)
	}
}
