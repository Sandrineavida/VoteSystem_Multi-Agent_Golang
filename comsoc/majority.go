package comsoc

import (
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

func MajoritySWF(p types.Profile) (types.Count, error) {
	err := checkProfile(p)
	if err != nil {
		return nil, err
	}

	count := make(types.Count)
	for _, voter := range p {
		count[voter[0]]++
	}

	return count, nil
}

func MajoritySCF(p types.Profile) ([]types.Alternative, error) {
	count, err := MajoritySWF(p)
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
