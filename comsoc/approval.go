package comsoc

import (
	"errors"
	"fmt"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

// Approval
// Le seuil est un tableau car chaque électeur peut avoir son propre seuil. C'est-à-dire que différents électeurs peuvent avoir des critères ou des nombres différents de candidats qu'ils approuvent. En attribuant un seuil spécifique à chaque électeur, le système de vote peut être rendu plus flexible et personnalisé.
//
// Par exemple, considérons trois électeurs et quatre candidats. Chaque électeur pourrait approuver un nombre différent de candidats basé sur ses propres jugements et critères. Par exemple :
// L'électeur 1 pourrait seulement approuver son premier choix
// L'électeur 2 pourrait approuver ses deux premiers choix
// L'électeur 3 pourrait approuver tous les candidats
// Ainsi, chaque électeur peut avoir son propre seuil, par exemple [1, 2, 4]. Ce tableau de seuils signifie que :
//
// Le seuil de l'électeur 1 est 1, donc il ne vote que pour son premier choix
// Le seuil de l'électeur 2 est 2, donc il vote pour ses deux premiers choix
// Le seuil de l'électeur 3 est 4, donc il vote pour tous les candidats

func ApprovalSWF(p types.Profile, thresholds []int) (types.Count, error) {
	if len(p) != len(thresholds) {
		return nil, errors.New("le nombre d'électeurs ne correspond pas à la longueur du tableau de seuils")
	}

	err := checkProfile(p)
	if err != nil {
		return nil, err
	}

	count := make(types.Count)
	for i, voter := range p {
		threshold := thresholds[i]
		if threshold < 0 || threshold >= len(voter) {
			return nil, fmt.Errorf("threshold %d pour le voter %d : mauvaise longueur ", threshold, i)
		}

		for j, alt := range voter {
			if j <= threshold {
				count[alt]++
			}
		}
	}

	return count, nil
}

func ApprovalSCF(p types.Profile, thresholds []int) ([]types.Alternative, error) {
	count, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}

	max := 0
	var bestAlts []types.Alternative
	for alt, votes := range count {
		if votes > max {
			max = votes
			bestAlts = []types.Alternative{alt}
		} else if votes == max {
			bestAlts = append(bestAlts, alt)
		}
	}

	return bestAlts, nil
}
