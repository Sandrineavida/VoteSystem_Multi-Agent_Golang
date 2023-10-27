package comsoc

import (
	"errors"
	"fmt"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

func CondorcetWinner(p types.Profile) (bestAlts []types.Alternative, err error) {
	err = checkProfile(p)
	if err != nil {
		return nil, err
	}

	win_time := make(map[types.Alternative]int) // Number of wins for each alternative

	for i := 1; i <= len(p[0])-1; i++ {
		for j := i + 1; j <= len(p[0]); j++ {
			count_i, count_j := 0, 0
			// Parcourir chaque préférence d'individu dans Profile pour compter les votes
			for _, prefs := range p {
				if isPref(types.Alternative(i), types.Alternative(j), prefs) {
					count_i++
				} else {
					count_j++
				}
			}
			// Resultat du vote : {i,j}, le gagnant de l'alt son nombre de victoires +1
			if count_i > count_j {
				win_time[types.Alternative(i)]++
			} else {
				win_time[types.Alternative(j)]++
			}
		}
	}

	fmt.Println(win_time)

	// Trouver le gagnant Condorcet, son win_time[Alternative(i)] devrait être le nombre total de candidats -1
	for key, value := range win_time {
		if value == len(p[0])-1 {
			bestAlts = append(bestAlts, types.Alternative(key))
			return bestAlts, nil
		}
	}

	return nil, errors.New("pas de vainqueur Condorcet")
}
