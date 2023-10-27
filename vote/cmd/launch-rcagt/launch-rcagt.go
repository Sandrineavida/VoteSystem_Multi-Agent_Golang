package main

import (
	"fmt"
	"math/rand"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/vote/restclientagent"
)

func main() {
	nCandidat := 10
	permutation := rand.Perm(nCandidat)
	prefs := make([]types.Alternative, 0, 10)
	for i := range permutation {
		permutation[i] += 1
		prefs[i] = types.Alternative(permutation[i])
	}
	ag := restclientagent.NewRestClientAgent("id1", "http://localhost:8000", prefs, rand.Intn(10), nCandidat)
	ag.Start()
	fmt.Scanln()
}
