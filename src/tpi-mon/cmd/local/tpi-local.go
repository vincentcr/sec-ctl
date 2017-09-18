package main

import (
	"log"
	"tpi-mon/api"
	"tpi-mon/tpi"
)

func main() {
	c, err := tpi.NewLocalClient("127.0.0.1", 9751, "aBcDe1")
	if err != nil {
		log.Panicln(err)
	}

	api.Run(c, 9750)
}
