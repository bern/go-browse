package utils

import "github.com/bern/go-browse/cmd/go-browse/models"

// BySpecificity sorts a slice of selectors with the highest specificity first
type BySpecificity []models.Selector

func (a BySpecificity) Len() int {
	return len(a)
}

func (a BySpecificity) Less(i, j int) bool {
	specI := CalculateSpecificity(a[i])
	specJ := CalculateSpecificity(a[j])

	if specI.IDSpecificity > specJ.IDSpecificity {
		return true
	} else if specI.IDSpecificity < specJ.IDSpecificity {
		return false
	}

	if specI.ClassSpecificity > specJ.ClassSpecificity {
		return true
	} else if specI.ClassSpecificity < specJ.ClassSpecificity {
		return false
	}

	if specI.ElementSpecificity > specJ.ElementSpecificity {
		return true
	} else if specI.ElementSpecificity < specJ.ElementSpecificity {
		return false
	}

	// must be sorted by source order
	return false
}

func (a BySpecificity) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ByMatchedRuleSpecificityAscending sorts a slice of rules with the LOWEST specificity first
type ByMatchedRuleSpecificityAscending []models.MatchedRule

func (a ByMatchedRuleSpecificityAscending) Len() int {
	return len(a)
}

func (a ByMatchedRuleSpecificityAscending) Less(i, j int) bool {
	specI := a[i].Specificity
	specJ := a[j].Specificity

	if specI.IDSpecificity < specJ.IDSpecificity {
		return true
	} else if specI.IDSpecificity > specJ.IDSpecificity {
		return false
	}

	if specI.ClassSpecificity < specJ.ClassSpecificity {
		return true
	} else if specI.ClassSpecificity > specJ.ClassSpecificity {
		return false
	}

	if specI.ElementSpecificity < specJ.ElementSpecificity {
		return true
	} else if specI.ElementSpecificity > specJ.ElementSpecificity {
		return false
	}

	// must be sorted by source order
	return false
}

func (a ByMatchedRuleSpecificityAscending) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// CalculateSpecificity returns the Specificity block for a given selector
func CalculateSpecificity(s models.Selector) models.Specificity {
	idSpecificity := 0
	classSpecificity := 0
	elementSpecificity := 0

	if s.ID != nil {
		idSpecificity = 1
	}

	if s.Classes != nil {
		classSpecificity = len(*s.Classes)
	}

	if s.TagName != nil {
		elementSpecificity = 1
	}

	return models.Specificity{
		IDSpecificity:      idSpecificity,
		ClassSpecificity:   classSpecificity,
		ElementSpecificity: elementSpecificity,
	}
}
