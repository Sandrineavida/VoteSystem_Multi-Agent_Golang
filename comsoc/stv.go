package comsoc

import (
	"errors"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

func STV_SWF(p types.Profile) (types.Count, error) {
	err := checkProfile(p)
	if err != nil {
		return nil, err
	}

	counts := make(types.Count)
	for _, voter := range p {
		if len(voter) > 0 {
			counts[voter[0]]++
		}
	}
	return counts, nil
}

func STV_SCF(p types.Profile) (bestAlts []types.Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	counts := make(types.Count) // Initialisation de tous les comptes à 0
	totalVotes := len(p)
	threshold := totalVotes/2 + 1
	maxIterations := len(p[0]) - 1 // Supposer que chaque personne a voté un vote complet, donc le nombre total de candidats est la longueur de p[0]

	for iteration := 0; iteration < maxIterations; iteration++ {
		// 计数
		for _, voter := range p {
			if len(voter) > 0 {
				counts[voter[0]]++
				if counts[voter[0]] >= threshold {
					return []types.Alternative{voter[0]}, nil // Atteindre le seuil, sélectionnez le gagnant
				}
			}
		}

		// Trouver le nombre de votes le plus bas
		var minCount = totalVotes + 1 // initialisation à une valeur plus grande que le nombre maximum de votes possible
		for _, count := range counts {
			if count >= 0 && count < minCount {
				minCount = count
			}
		}

		// Eliminer les candidats avec le nombre de votes le plus bas
		eliminated := make(map[types.Alternative]bool)
		for alt, count := range counts {
			if count == minCount {
				eliminated[alt] = true
			}
		}

		// Si tous les candidats restants ont le même nombre de votes, ils sont tous gagnants
		if len(eliminated) == len(counts) || iteration == maxIterations-1 {
			maxCount := 0
			for _, count := range counts {
				if count > maxCount {
					maxCount = count
				}
			}

			for alt, count := range counts {
				if count == maxCount {
					bestAlts = append(bestAlts, alt)
				}
			}
			return bestAlts, nil
		}

		// Sinon, transférer les votes des candidats éliminés
		for i, voter := range p {
			// Eliminer les candidats éliminés
			newPrefs := []types.Alternative{}
			for _, alt := range voter {
				if !eliminated[alt] {
					newPrefs = append(newPrefs, alt)
				}
			}
			p[i] = newPrefs
		}

		// Réinitialiser les comptes
		for alt := range counts {
			counts[alt] = 0
		}
	}

	// Si on arrive ici, cela signifie qu'aucun gagnant n'a été sélectionné
	return nil, errors.New("no winner selected after maximum iterations")
}
