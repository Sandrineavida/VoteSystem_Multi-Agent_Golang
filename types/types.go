package types

import (
	"errors"
	"time"
)

type Alternative int           // type pour les alternatives (candidats)
type Profile [][]Alternative   // type pour les profils (liste des préférences des électeurs)
type Count map[Alternative]int // type pour les résultats (map des scores(value) des candidats(key))

// Ballot
type Ballot struct {
	ID         string          // ID unique d'un ballot
	Rule       string          // Règle de vote
	Deadline   time.Time       // Date limite de vote
	Votes      [][]Alternative // Tous les votes (Profile)
	HasCounted bool            // Indique si le dépouillement a été effectué. Une fois la date limite dépassée et les résultats calculés, il est défini sur true.
	VoterIDs   []string        // Liste des ID des votants ayant déjà voté
	Alts       int             // Nombre de candidats
	TieBreak   []Alternative   // L'ordre prédéfini des candidats utilisé quand on a un "tie"
	Result     Alternative     // Résultat du vote (Un gagnant)
	Thresholds []int           // Dans le cadre d'Approval, le nombre de votes pour chaque électeur
	Seuil      int             // Seuil dans le cadre d'Approval (qui doit <= Alts int)
}

// VoteRequest
type VoteRequest struct {
	AgentID  string        `json:"agent-id"`  // ID de l'électeur
	BallotID string        `json:"ballot-id"` // ID du ballot
	Prefs    []Alternative `json:"prefs"`     // Préférences de l'électeur
	Options  int           `json:"options"`   // Nombre d'options considerées (Utilisé uniquement dans le cadre d'Approval)
}

// NewBallotRequest
type NewBallotRequest struct {
	Rule     string        `json:"rule"`            // Règle de vote (ex.: "majority", "borda", etc.)
	Deadline time.Time     `json:"deadline"`        // ate limite pour voter, format RFC 3339 (ex.: "2023-10-09T23:05:08+02:00")
	VoterIDs []string      `json:"voter-ids"`       // Liste des IDs des votants (ex.: ["ag_id1", "ag_id2", "ag_id3"])
	Alts     int           `json:"#alts"`           // Nombre de candidats (ex.: 12)
	TieBreak []Alternative `json:"tie-break"`       // L'ordre prédéfini des candidats utilisé quand on a un "tie" (ex.: [4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10])
	Seuil    int           `json:"seuil,omitempty"` // Seuil dans le cadre d'Approval (qui doit <= Alts int) (ex.: 10)
}

// ResultVoteRequest est la requête pour la méthode /result
type ResultVoteRequest struct {
	BallotID string `json:"ballot-id"` // ID du ballot
}

// Type de base pour les réponses, contenant un indicateur de succès et un message éventuel.
type BaseResponse struct {
	Success bool   `json:"success"` // Indique si la demande a réussi
	Message string `json:"message"` // Message d'erreur ou de succès potentiel
}

// NewBallotResponse est la réponse à la requête /new_ballot
type NewBallotResponse struct {
	BaseResponse        // Intégration de la structure de réponse de base
	BallotID     string `json:"ballot_id,omitempty"` // ID du bulletin nouvellement créé
}

// VoteResponse est la réponse à la requête /vote
type VoteResponse struct {
	BaseResponse
}

// ResultVoteResponse est la réponse à la requête /result
type ResultVoteResponse struct {
	BaseResponse
	Winner Alternative `json:"winner"`
	// Ranking []int       `json:"ranking,omitempty"` // optionnel
}

// Verifier si la requête de creer un nouveau ballot est valide
func (req *NewBallotRequest) Validate_NewBallotRequest() error {
	if req.Rule == "" {
		return errors.New("missing voting rule")
	}
	if req.Deadline.IsZero() {
		return errors.New("missing deadline")
	}
	if req.Alts <= 0 {
		return errors.New("invalid number of alternatives")
	}
	//if len(req.TieBreak) != req.Alts {
	//	return fmt.Errorf("tie break length (%d) does not match number of alternatives (%d)", len(req.TieBreak), req.Alts)
	//}
	return nil
}
