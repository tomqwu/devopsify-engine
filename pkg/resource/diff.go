package resource

import "fmt"

// Diff compares desired resources against actual resources and returns drift results.
// It detects missing, added, and modified resources by matching on ID.
func Diff(desired, actual []Resource) []DriftResult {
	var results []DriftResult

	desiredMap := make(map[string]Resource, len(desired))
	for _, r := range desired {
		desiredMap[r.ID] = r
	}

	actualMap := make(map[string]Resource, len(actual))
	for _, r := range actual {
		actualMap[r.ID] = r
	}

	// Check for missing or modified resources
	for id, d := range desiredMap {
		a, exists := actualMap[id]
		if !exists {
			results = append(results, DriftResult{
				ResourceID:   id,
				ResourceType: d.Type,
				DriftType:    DriftTypeMissing,
				Message:      fmt.Sprintf("resource %s (%s) exists in desired state but not in actual", d.Name, id),
			})
			continue
		}

		if d.State != "" && d.State != a.State {
			results = append(results, DriftResult{
				ResourceID:   id,
				ResourceType: d.Type,
				DriftType:    DriftTypeModified,
				Field:        "state",
				Expected:     string(d.State),
				Actual:       string(a.State),
				Message:      fmt.Sprintf("resource %s state: expected %s, got %s", d.Name, d.State, a.State),
			})
		}

		// Check tag drift
		for key, val := range d.Tags {
			if actualVal, ok := a.Tags[key]; !ok {
				results = append(results, DriftResult{
					ResourceID:   id,
					ResourceType: d.Type,
					DriftType:    DriftTypeModified,
					Field:        fmt.Sprintf("tags.%s", key),
					Expected:     val,
					Actual:       "",
					Message:      fmt.Sprintf("resource %s missing tag %s=%s", d.Name, key, val),
				})
			} else if actualVal != val {
				results = append(results, DriftResult{
					ResourceID:   id,
					ResourceType: d.Type,
					DriftType:    DriftTypeModified,
					Field:        fmt.Sprintf("tags.%s", key),
					Expected:     val,
					Actual:       actualVal,
					Message:      fmt.Sprintf("resource %s tag %s: expected %s, got %s", d.Name, key, val, actualVal),
				})
			}
		}
	}

	// Check for added resources (in actual but not desired)
	for id, a := range actualMap {
		if _, exists := desiredMap[id]; !exists {
			results = append(results, DriftResult{
				ResourceID:   id,
				ResourceType: a.Type,
				DriftType:    DriftTypeAdded,
				Message:      fmt.Sprintf("resource %s (%s) exists in actual state but not in desired", a.Name, id),
			})
		}
	}

	return results
}
