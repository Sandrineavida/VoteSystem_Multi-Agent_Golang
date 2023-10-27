package restclientagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

// Definire la structure d'un agent client qui est soit électeur soit agent qui génère des ballots
type RestClientAgent struct {
	id        string              // ID de l'agent
	url       string              // URL du serveur
	prefs     []types.Alternative // Préférences de l'agent (Si c'est un électeur)
	seuil     int                 // Valabe uniquement pour Approval (Si c'est un électeur)
	options   int                 // Nombre d'options considerées, valabe uniquement pour Approval (Si c'est un électeur) (N.B.: defini apres dans la fonction doVoteRequest() en fonction de seuil)
	nCandidat int                 // Nombre de candidats
}

// Instancier un agent client (type RestClientAgent)
func NewRestClientAgent(id string, url string, prefs []types.Alternative, seuil int, nCanditat int) *RestClientAgent {
	return &RestClientAgent{id, url, prefs, seuil, -1, nCanditat}
}

// Cote client: traiter la reponse de serveur pour la requete de voter (Cote serveur: doVoteRequest)
func (rca *RestClientAgent) treatVoteResponse(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var resp types.VoteResponse
	json.Unmarshal(buf.Bytes(), &resp)

	if resp.Success {
		return "Vote successful"
	} else {
		return fmt.Sprintf("Vote failed: %s", resp.Message)
	}
}

// Cote client: traiter la reponse de serveur pour la requete de creer un nouveau ballot (Cote serveur: doNewBallotRequest)
func (rca *RestClientAgent) treatNewBallotResponse(r *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var resp types.NewBallotResponse
	json.Unmarshal(buf.Bytes(), &resp)

	if resp.Success {
		return fmt.Sprintf("New ballot created successfully, Ballot ID: %s", resp.BallotID)
	} else {
		return fmt.Sprintf("Failed to create new ballot: %s", resp.Message)
	}
}

// Cote client: traiter la reponse de serveur pour la requete de resultat (Cote serveur: doResultVoteRequest)
func (rca *RestClientAgent) treatResultVoteResponse(r *http.Response) (types.Alternative, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var resp types.ResultVoteResponse
	err := json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf(resp.Message)
	}

	return resp.Winner, nil
}

// Cote client: traiter la reponse de serveur pour la requete de compter le nombre de requetes (Cote serveur: doCountRequest)
func (rca *RestClientAgent) treatCountResponse(r *http.Response) (int, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var resp struct {
		Count   int    `json:"count"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	err := json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("failed to retrieve count: %s", resp.Message)
	}

	return resp.Count, nil
}

// Cote client: demander de créer un nouveau ballot
func (rca *RestClientAgent) doNewBallotRequest(rule string, deadline string, voterIDs []string, alts int, tieBreak []types.Alternative) (string, error) {
	// transformer deadline en type time.Time
	deadlineTime, err := time.Parse(time.RFC3339, deadline)
	if err != nil {
		return "", err
	}

	// Affecter les valeurs aux champs de la structure NewBallotRequest et stocker dans la variable req
	req := types.NewBallotRequest{
		Rule:     rule,
		Deadline: deadlineTime,
		VoterIDs: voterIDs,
		Alts:     alts,
		TieBreak: tieBreak,
	}

	// Si le client souhaite creer un ballot de type "approval", il faut affecter le champ "Seuil" de la structure NewBallotRequest
	if rule == "approval" {
		req.Seuil = rca.seuil
	}

	// Valider la nouvelle demande de vote
	if err := req.Validate_NewBallotRequest(); err != nil {
		return "", err
	}

	// Transforme la requete en format JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// HTTP POST: Envoie la requete au serveur et recupere la reponse dans la variable "response"
	response, err := http.Post(rca.url+"/new_ballot", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", response.StatusCode, response.Status)
		return "", err
	}

	// Traitement de la reponse du serveur et retourne le resultat
	return rca.treatNewBallotResponse(response), nil
}

// Cote client: demander de voter
func (rca *RestClientAgent) doVoteRequest(ballotID string) (string, error) {
	// Creer une requete de type "types.VoteRequest"
	req := types.VoteRequest{
		AgentID:  rca.id,
		BallotID: ballotID,
		Prefs:    rca.prefs,
	}

	if ballotID == "vote_approval" { // Si on est dans le contexte d'Approval
		rand.Seed(time.Now().UnixNano())         // Initialisation de la graine aléatoire
		rca.options = rand.Intn(rca.seuil-1) + 1 // Génération aléatoire de la valeur Options (entre 1 et seuil-1)
		req.Options = rca.options                // Mise à jour du champ Options de la requete
	}

	// Transforme la requete en format JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// Envoie la requete pour voter au serveur et recupere la reponse dans la variable "response"
	response, err := http.Post(rca.url+"/vote", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Si le code de statut n'est pas 200 OK, on retourne une erreur
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get count: %d %s", response.StatusCode, response.Status)
	}

	// Traitement de la reponse du serveur et retourne le resultat
	return rca.treatVoteResponse(response), nil
}

// Cote client: demander le resultat d'un ballot
func (rca *RestClientAgent) doResultVoteRequest(ballotID string) (winner types.Alternative, err error) {
	// Creer une requete de type "types.ResultVoteRequest"
	req := types.ResultVoteRequest{
		BallotID: ballotID,
	}
	// HTTP POST: Envoie l'ID du ballot a "/result" et recupere la reponse dans la variable "response"
	jsonData, _ := json.Marshal(req)
	response, err := http.Post(rca.url+"/result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return types.Alternative(0), err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get count: %d %s", response.StatusCode, response.Status)
	}

	// Traitement de la reponse du serveur
	winner, err = rca.treatResultVoteResponse(response)
	if err != nil {
		return types.Alternative(0), err
	}

	return winner, nil
}

// Cote client: demander le nombre de requetes
func (rca *RestClientAgent) doCountRequest() (int, error) {
	// 发送 HTTP GET 请求到服务器的 /count 端点

	// HTTP GET: Recupere la reponse dans la variable "response"
	response, err := http.Get(rca.url + "/count")
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get count: %d %s", response.StatusCode, response.Status)
	}

	// Traitement de la reponse du serveur et retourne le resultat
	return rca.treatCountResponse(response)
}

// Fonc. de Start pour l'agent qui génère les ballots
func (rca *RestClientAgent) CreateBallotStart() {
	log.Println("Creating different ballots...")

	// Definir un tableau "ballots" qui contient le nom des differentes regles de vote
	ballots := []string{"borda", "majority", "copeland", "approval", "stv"} //! 5 ballots => nb. de requetes += 5

	// Pour chaque regle de vote, creer un nouveau ballot
	for _, ballot := range ballots {
		// tieBreak : l'ordre prédéfini des candidats utilisé quand on a un "tie"
		// ex.: si nCandidat = 10, tieBreak = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		tieBreak := make([]types.Alternative, rca.nCandidat)
		for i := 0; i < rca.nCandidat; i++ {
			tieBreak[i] = types.Alternative(i + 1)
		}

		// Demander au serveur de creer un nouveau ballot et recuperer la reponse dans la variable "response"
		response, err := rca.doNewBallotRequest(ballot, time.Now().Add(10*time.Second).Format(time.RFC3339), nil, rca.nCandidat, tieBreak)
		if err != nil {
			log.Fatalf("[%s] Error: %s\n", rca.id, err.Error())
			continue
		}
		log.Printf("[POST][%s] Request to create ballot, id = %v, %s\n", rca.id, ballot, response)
	}
}

// Fonc. de Start pour l'agent qui est électeur
func (rca *RestClientAgent) Start() {
	log.Printf("démarrage de %s", rca.id)

	ballots := []string{"vote_borda", "vote_majority", "vote_approval", "vote_stv", "vote_copeland"}

	// Voter pour chaque ballot
	for _, ballot := range ballots { //! nb. de requetes += 5 * nAgent
		res, err := rca.doVoteRequest(ballot)
		if err != nil {
			log.Fatal(rca.id, "error:", err.Error())
		} else {
			if ballot == "vote_approval" {
				log.Printf("[POST][%s] voting to %s, preferences = %v, %s\n", rca.id, ballot, rca.prefs[:rca.options], res)
			} else {
				log.Printf("[POST][%s] voting to %s, preferences = %v, %s\n", rca.id, ballot, rca.prefs, res)
			}
		}
	}

	time.Sleep(10 * time.Second)

	// Demander le resultat de chaque ballot
	for _, ballot := range ballots { //! nb. de requetes += 5 * nAgent
		res, err := rca.doResultVoteRequest(ballot)
		if err != nil {
			log.Fatal(rca.id, "error:", err.Error())
		} else {
			log.Printf("[GET][%s] requesting result of election [%s], Alternative elected = %v\n", rca.id, ballot, res)
		}
	}

	time.Sleep(10 * time.Second)

	// Demander le nombre de requetes
	res3, err3 := rca.doCountRequest()
	if err3 != nil {
		log.Fatal(rca.id, "error:", err3.Error())
	} else {
		log.Printf("[GET][%s] requesting server request count, count = %v\n", rca.id, res3) //! nb. de requetes doit = 5 + 5 * nAgent * 2
	}
}
