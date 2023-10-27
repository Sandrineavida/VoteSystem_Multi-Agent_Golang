package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/vote/restclientagent"
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/vote/restserveragent"
)

func main() {

	// creer un fichier pour stocker les logs
	// file, err := os.Create("./td5/result3.txt")
	file, err := os.Create("./result3.txt ")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	log.SetOutput(file)

	const nAgent = 100
	const nCandidat = 16
	const url1 = ":8080"
	const url2 = "http://localhost:8080"
	const seuil = 10 // pour Approval

	ballotCreationAgt := restclientagent.NewRestClientAgent("ballot_creation_agt", url2, nil, seuil, nCandidat) // client qui génère les ballots
	clAgts := make([]restclientagent.RestClientAgent, 0, nAgent)                                                // tableau des agents clients qui sont des électeurs
	servAgt := restserveragent.NewRestServerAgent(url1)                                                         // agent serveur

	log.Println("démarrage du serveur...")
	go servAgt.Start()
	log.Println("démarrage du client qui génère les ballots...")
	go ballotCreationAgt.CreateBallotStart()

	time.Sleep(5 * time.Second)

	log.Println("démarrage des clients qui sont des électeurs...")
	// Instancier les agents clients qui sont des électeurs
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nAgent; i++ {
		// generer une preference aleatoire pour chaque electeur
		permutation := rand.Perm(nCandidat)
		prefs := make([]types.Alternative, nCandidat)
		for i := range permutation {
			permutation[i] += 1
			prefs[i] = types.Alternative(permutation[i])
		}
		// generer un id pour chaque electeur
		id := fmt.Sprintf("id%02d", i)
		// instancier un agent client (electeur)
		agt := restclientagent.NewRestClientAgent(id, url2, prefs, seuil, nCandidat)
		// ajouter l'agent client dans le tableau
		clAgts = append(clAgts, *agt)
	}
	// Démarrer les agents clients qui sont des électeurs
	for _, agt := range clAgts {
		// attention, obligation de passer par cette lambda pour faire capturer la valeur de l'itération par la goroutine
		func(agt restclientagent.RestClientAgent) {
			go agt.Start()
		}(agt)
	}

	fmt.Scanln()
}
