package comsoc

import "gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"

// La fonction d'usine accepte une fonction de bien-être social (SWF) ou une fonction de choix social (SCF) et une fonction de décision en cas d'égalité (TieBreak), et renvoie une nouvelle fonction. Cette nouvelle fonction appelle la SWF ou SCF donnée, puis utilise TieBreak pour résoudre toute égalité.
// SWFFactory accepte une fonction swf et une fonction tieBreak en tant que paramètres. Elle renvoie une nouvelle fonction qui appelle d'abord swf, puis tieBreak en cas de besoin pour résoudre une égalité. La nouvelle fonction renvoie une tranche contenant les options gagnantes et une erreur (si présente).
// SCFFactory utilise une fonction scf pour déterminer la meilleure option. En cas d'égalité, elle utilise la fonction tieBreak pour déterminer le gagnant final. La nouvelle fonction renvoie l'option gagnante et une erreur (si présente).

func SWFFactory(swf func(p types.Profile) (types.Count, error), tieBreak func([]types.Alternative) (types.Alternative, error)) func(types.Profile) ([]types.Alternative, error) {
	return func(p types.Profile) ([]types.Alternative, error) {
		count, err := swf(p)
		if err != nil {
			return nil, err
		}

		maxCount := maxCount(count)
		if len(maxCount) == 1 {
			return maxCount, nil
		}

		winner, err := tieBreak(maxCount)
		if err != nil {
			return nil, err
		}
		return []types.Alternative{winner}, nil
	}
}

func SCFFactory(scf func(p types.Profile) ([]types.Alternative, error), tieBreak func([]types.Alternative) (types.Alternative, error)) func(types.Profile) (types.Alternative, error) {
	return func(p types.Profile) (types.Alternative, error) {
		bestAlts, err := scf(p)
		if err != nil {
			return -1, err
		}

		if len(bestAlts) == 1 {
			return bestAlts[0], nil
		}

		return tieBreak(bestAlts)
	}
}
