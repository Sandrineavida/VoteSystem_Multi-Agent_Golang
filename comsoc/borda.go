package comsoc

import (
	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

// Borda
func BordaSWF(p types.Profile) (types.Count, error) {
	err := checkProfile(p)
	if err != nil {
		return nil, err
	}

	count := make(types.Count)
	for _, voter := range p {
		for i, alt := range voter {
			count[alt] += len(voter) - i - 1
		}
	}

	return count, nil
}

func BordaSCF(p types.Profile) ([]types.Alternative, error) {
	count, err := BordaSWF(p)
	if err != nil {
		return nil, err
	}

	max := -1
	var bestAlts []types.Alternative
	for alt, score := range count {
		if score > max {
			max = score
			bestAlts = []types.Alternative{alt}
		} else if score == max {
			bestAlts = append(bestAlts, alt)
		}
	}

	return bestAlts, nil
}
