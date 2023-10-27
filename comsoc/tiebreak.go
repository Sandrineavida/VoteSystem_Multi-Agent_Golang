package comsoc

import (
	"errors"
	"fmt"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

// Traiter une situation de "tie"
func TieBreakFactory(orderedAlts []types.Alternative) func([]types.Alternative) (types.Alternative, error) {
	return func(candidates []types.Alternative) (types.Alternative, error) {
		if len(candidates) == 0 {
			return 0, errors.New("la liste des candidats est vide")
		}

		// map utilisé pour trouver la position d'un candidat dans l'ordre prédéfini
		ranking := make(map[types.Alternative]int)
		for i, alt := range orderedAlts {
			ranking[alt] = i
		}

		// trouver le candidat avec le rang le plus bas
		winner := candidates[0]
		minRank := ranking[winner]
		for _, candidate := range candidates {
			if rank, found := ranking[candidate]; found {
				if rank < minRank {
					winner = candidate
					minRank = rank
				}
			} else {
				return 0, fmt.Errorf("le candidat %v n'est pas dans la liste de séquence prédéfinie", candidate)
			}
		}

		return winner, nil
	}
}
