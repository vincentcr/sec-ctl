package main

import "envisalink"

func main() {
	c := envisalink.NewClient()
	c.Connect("192.168.86.12", "Q4m1gh")

}
