package main

import (
	"log"
	"tpi-mon/mock"
)

func main() {
	if err := mock.Run("0.0.0.0", 9751, 9752, "mock-tpi-state.json"); err != nil {
		log.Panicln(err)
	}
}
