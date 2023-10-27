package main

import (
	"fmt"

	ras "gitlab.utc.fr/sunhudie/ia04-projet-par-binome/vote/restserveragent"
)

func main() {
	server := ras.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
