package comsoc

import (
	"errors"
	"fmt"

	"gitlab.utc.fr/sunhudie/ia04-projet-par-binome/types"
)

// renvoie l'indice ou se trouve alt dans prefs
func rank(alt types.Alternative, prefs []types.Alternative) int {
	for i, a := range prefs {
		if a == alt {
			return i
		}
	}
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func isPref(alt1, alt2 types.Alternative, prefs []types.Alternative) bool {
	return rank(alt1, prefs) < rank(alt2, prefs)
}

// renvoie les meilleures alternatives pour un décomtpe donné
func maxCount(count types.Count) (bestAlts []types.Alternative) {
	maxVal := -1
	for _, val := range count {
		if val > maxVal {
			maxVal = val
		}
	}

	for alt, val := range count {
		if val == maxVal {
			bestAlts = append(bestAlts, alt)
		}
	}
	return
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois par préférences
func checkProfile(prefs types.Profile) error {
	var errs []error
	for i, voter := range prefs {
		altSet := make(map[types.Alternative]bool)
		for _, alt := range voter {
			if _, exists := altSet[alt]; exists {
				errs = append(errs, fmt.Errorf("duplicate alternative in voter profile at index %d", i))
			}
			altSet[alt] = true
		}

		if len(altSet) != len(voter) {
			errs = append(errs, fmt.Errorf("incomplete voter profile at index %d", i))
		}
	}

	if len(errs) > 0 {
		var errStr string
		for _, err := range errs {
			errStr += err.Error() + "; "
		}
		return errors.New(errStr)
	}

	return nil
}

// Vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs types.Profile, alts []types.Alternative) error {
	err := checkProfile(prefs)
	if err != nil {
		return err // 如果 checkProfile 返回了错误，直接返回这个错误
	}

	var errs []error

	altSet := make(map[types.Alternative]bool)
	for _, alt := range alts {
		altSet[alt] = true
	}

	for i, voter := range prefs {
		voterAltSet := make(map[types.Alternative]bool)
		for _, alt := range voter {
			if !altSet[alt] {
				errs = append(errs, fmt.Errorf("unknown alternative in voter profile at index %d", i))
			}
			if _, exists := voterAltSet[alt]; exists {
				errs = append(errs, fmt.Errorf("duplicate alternative in voter profile at index %d", i))
			}
			voterAltSet[alt] = true
		}

		if len(voterAltSet) != len(alts) {
			errs = append(errs, fmt.Errorf("incomplete or extra alternatives in voter profile at index %d", i))
		}
	}

	if len(errs) > 0 {
		var errStr string
		for _, err := range errs {
			errStr += err.Error() + "; "
		}
		return errors.New(errStr)
	}

	return nil
}
