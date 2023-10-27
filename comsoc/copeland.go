package comsoc

import (
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

func CopelandSWF(p types.Profile) (types.Count, error) {
	err := checkProfile(p)
	if err != nil {
		return nil, err
	}

	// initialiser le compteur de Copeland
	copelandCount := make(types.Count)

	// Comparer chaque paire de candidats
	for _, voter := range p {
		for i := 0; i < len(voter); i++ {
			for j := i + 1; j < len(voter); j++ {
				if isPref(voter[i], voter[j], voter) {
					copelandCount[voter[i]]++
					copelandCount[voter[j]]--
				} else if isPref(voter[j], voter[i], voter) {
					copelandCount[voter[j]]++
					copelandCount[voter[i]]--
				}
			}
		}
	}

	return copelandCount, nil
}

func CopelandSCF(p types.Profile) ([]types.Alternative, error) {
	copelandCount, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}

	// Trouver le candidat avec le score le plus élevé
	var maxScore int
	for _, score := range copelandCount {
		if score > maxScore {
			maxScore = score
		}
	}

	// On peut avoir plusieurs candidats avec le même score le plus élevé
	var bestAlts []types.Alternative
	for alt, score := range copelandCount {
		if score == maxScore {
			bestAlts = append(bestAlts, alt)
		}
	}

	//if len(bestAlts) == 0 {
	//	return nil, errors.New("Aucun gagnant de Copeland trouvé")
	//}

	return bestAlts, nil
}
