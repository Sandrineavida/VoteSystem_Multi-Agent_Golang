package restserveragent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/comsoc"
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

// Definir la structure RestServerAgent qui contient un mutex, un id, un compteur de requêtes, un map de ballots et une adresse.
type RestServerAgent struct {
	sync.Mutex
	id       string                  // id d'un agent de serveur
	reqCount int                     // compteur de requêtes
	ballots  map[string]types.Ballot // map de ballots
	addr     string                  // adresse
}

// Instancier un agent de serveur (type RestServerAgent)
func NewRestServerAgent(addr string) *RestServerAgent {
	return &RestServerAgent{
		id:      addr,
		addr:    addr,
		ballots: make(map[string]types.Ballot),
	}
}

// Test de la méthode
func (rsa *RestServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405 Method Not Allowed
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

// Decoder les requetes de type NewBallotRequest
func (*RestServerAgent) decodeNewBallotRequest(r *http.Request) (req types.NewBallotRequest, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		return req, err // Erreur de lecture du body
	}
	// Lire la requete et ecrire l'objet dans req envoyé en parametre
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		return req, err // Erreur de décodage JSON
	}
	return req, nil
}

// Decoder les requetes de type VoteRequest
func (*RestServerAgent) decodeVoteRequest(r *http.Request) (req types.VoteRequest, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		return req, err // Erreur de lecture du body
	}
	// Lire la requete et ecrire l'objet dans req envoyé en parametre
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		return req, err // Erreur de décodage JSON
	}
	return req, nil
}

// Decoder les requetes de type ResultVoteRequest
func (*RestServerAgent) decodeResultVoteRequest(r *http.Request) (req types.ResultVoteRequest, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		return req, err // Erreur de lecture du body
	}
	// Lire la requete et ecrire l'objet dans req envoyé en parametre
	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		return req, err // Erreur de décodage JSON
	}
	return req, nil
}

// Cote serveur: traiter la requete pour creer un nouveau ballot
func (rsa *RestServerAgent) doReqNewBallot(w http.ResponseWriter, r *http.Request) {

	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++ //! Incrementer le compteur de requetes

	if !rsa.checkMethod("POST", w, r) {
		return
	}

	req, err := rsa.decodeNewBallotRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprint(w, err.Error())
		return
	}

	// Verifier si la requête de creer un nouveau ballot est valide
	err = req.Validate_NewBallotRequest()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprint(w, err.Error())
		return
	}

	// Creer un nouveau ballot
	ballotID := "vote_" + req.Rule // Creer un ID pour le nouveau ballot (ex.: vote_majority, vote_borda, etc.)
	ballot := types.Ballot{
		ID:         ballotID,
		Rule:       req.Rule,
		Deadline:   req.Deadline,
		VoterIDs:   req.VoterIDs,
		Alts:       req.Alts,
		TieBreak:   req.TieBreak,
		Seuil:      req.Seuil,
		HasCounted: false,
	}

	// Stocker le nouveau ballot dans le serveur
	rsa.ballots[ballotID] = ballot

	// Creer une reponse correspondante qui indique que le ballot a été créé avec succès et contient l'ID du nouveau ballot
	resp := types.NewBallotResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
			Message: "Ballot created successfully",
		},
		BallotID: ballotID,
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
	serial, _ := json.Marshal(resp)   // Convertir la reponse en JSON
	w.Write(serial)                   // Ecrire la reponse dans le body
}

// Cote serceur: traiter la requete pour voter
func (rsa *RestServerAgent) doVote(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++ //! Incrementer le compteur de requetes

	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// decoder la requete pour voter
	req, err := rsa.decodeVoteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprint(w, "Invalid vote request: ", err.Error())
		return
	}

	// Verifier si le ballot d'ID req.BallotID existe dans le serveur
	ballot, exists := rsa.ballots[req.BallotID]
	if !exists {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		fmt.Fprint(w, "Invalid ballot ID")
		return
	}

	// Verifier si le ballot est encore ouvert
	if time.Now().After(ballot.Deadline) {
		w.WriteHeader(http.StatusServiceUnavailable) // 503 Service Unavailable
		fmt.Fprint(w, "Voting deadline has passed for this ballot")
		return
	}

	// Verifier si l'agent a déjà voté pour CE ballot
	hasVoted := false
	for _, voteID := range ballot.VoterIDs {
		if voteID == req.AgentID {
			hasVoted = true
			break
		}
	}
	if hasVoted {
		w.WriteHeader(http.StatusForbidden) // 403 Forbidden vote already cast
		fmt.Fprint(w, "Vote already cast by this voter")
		return
	}

	// Pas de probleme; Mise à jour du ballot
	ballot.Votes = append(ballot.Votes, req.Prefs)             // Ajouter les préférences de l'agent dans le tableau de Votes (Profile) du ballot
	ballot.VoterIDs = append(ballot.VoterIDs, req.AgentID)     // Ajouter l'ID de l'agent dans le tableau de VoterIDs du ballot
	ballot.Thresholds = append(ballot.Thresholds, req.Options) // Ajouter Options, nb. valide de candidats préférés par l'agent, dans le tableau de Thresholds du ballot (Utilisé uniquement dans le cadre d'Approval)
	rsa.ballots[req.BallotID] = ballot

	// Creer une reponse correspondante qui indique que le vote a été accepté
	resp := types.VoteResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
			Message: "Vote accepted",
		},
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

// Donner le nombre de requêtes reçues par le serveur
func (rsa *RestServerAgent) doReqcount(w http.ResponseWriter, r *http.Request) {

	if !rsa.checkMethod("GET", w, r) {
		return
	}

	resp := struct {
		Success bool `json:"success"`
		Count   int  `json:"count"`
	}{
		Success: true,
		Count:   rsa.reqCount,
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	rsa.Lock()
	defer rsa.Unlock()
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

// Cote serveur: traiter la requete pour obtenir le résultat du vote
func (rsa *RestServerAgent) doReqresult(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()
	rsa.reqCount++ //! Incrementer le compteur de requetes

	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// decoder la requete pour obtenir le résultat du vote
	req, err := rsa.decodeResultVoteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// Verifier si le ballot d'ID req.BallotID existe dans le serveur
	ballot, exists := rsa.ballots[req.BallotID]
	if !exists {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
		fmt.Fprintf(w, "Ballot with ID %s does not exist", req.BallotID)
		return
	}

	// Verifier si le ballot est encore ouvert
	if time.Now().Before(ballot.Deadline) {
		w.WriteHeader(http.StatusTooEarly) // hotfix1015: 425 Too Early
		fmt.Fprintf(w, "Voting is still ongoing, results will be available after %v", ballot.Deadline)
		return
	}

	winner := ballot.Result
	// ranking := []int{2, 1, 4, 3} // optinnel

	// Creer une reponse correspondante qui indique que le résultat du vote a été obtenu avec succès
	resp := types.ResultVoteResponse{
		BaseResponse: types.BaseResponse{
			Success: true,
			Message: "Voting results retrieved successfully",
		},
		Winner: winner,
		// Ranking: ranking, // optionnel
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

// Essayer de Calculer le résultat du vote (toutes les 10 secondes)
func (rsa *RestServerAgent) calculateResultCheck() {
	ticker := time.NewTicker(10 * time.Second) // Envoi d'un signal toutes les 10 secondes
	defer ticker.Stop()

	for { // Boucle infinie; le calcul du résultat est effectué toutes les 10 secondes
		select {
		case <-ticker.C:
			rsa.Lock()
			for id, ballot := range rsa.ballots {
				if !ballot.HasCounted && time.Now().After(ballot.Deadline) {
					// Si le dépouillement n'a pas encore été effectué et que la date limite est dépassée pour ce ballot
					var err error
					var tempResult []types.Alternative // On peut avoir plusieurs gagnants (un tie)
					switch ballot.Rule {
					case "borda":
						tempResult, err = comsoc.BordaSCF(ballot.Votes)
					case "copeland":
						tempResult, err = comsoc.CopelandSCF(ballot.Votes)
					case "majority":
						tempResult, err = comsoc.MajoritySCF(ballot.Votes)
					case "approval":
						tempResult, err = comsoc.ApprovalSCF(ballot.Votes, ballot.Thresholds)
					case "stv":
						tempResult, err = comsoc.STV_SCF(ballot.Votes)
					default:
						log.Printf("Unknown voting rule: %s\n", ballot.Rule)
						err = errors.New("unknown voting rule")
					}

					if err != nil {
						log.Printf("Error calculating result for ballot %s: %v\n", id, err)
						continue
					}

					if len(tempResult) > 1 {
						// Si on a plusieurs gagnant, c-a-d on a un tie, on va utiliser TieBreakFactory pour créer une fonction de tieBreak
						tieBreak := comsoc.TieBreakFactory(ballot.TieBreak) // Passer l'ordre prédéfini des candidats comme paramètre

						winner, err := tieBreak(tempResult)
						if err != nil {
							log.Printf("Error breaking tie: %v\n", err)
							continue
						}

						// On note le gagnant de ce ballot dans le champ Result du ballot
						ballot.Result = winner

					} else { // Si on a un seul gagnant
						ballot.Result = tempResult[0] // On note directement le gagnant de ce ballot dans le champ Result du ballot
					}
					ballot.HasCounted = true // On note que le dépouillement a été effectué
					rsa.ballots[id] = ballot // Mis à jour des infos dans le serveur
					log.Printf("Ballot %s is closed and has been counted. Winner: %v\n", id, ballot.Result)
				}
			}
			rsa.Unlock()
		}
	}
}

// Lancer le serveur
func (rsa *RestServerAgent) Start() {
	rsa.ballots = make(map[string]types.Ballot)
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/vote", rsa.doVote)
	mux.HandleFunc("/count", rsa.doReqcount)
	mux.HandleFunc("/result", rsa.doReqresult)
	mux.HandleFunc("/new_ballot", rsa.doReqNewBallot)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go rsa.calculateResultCheck()
	go log.Fatal(s.ListenAndServe())
}
