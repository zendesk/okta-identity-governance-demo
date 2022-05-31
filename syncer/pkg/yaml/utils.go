package yaml

import (
	"reflect"
	"sort"
)

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// mapsEqual checks that dictionary of strings are equal regardless of sort
func mapsEqual(a, b map[string][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !sortedArrayEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// sortedArrayEqual checks that the values in slices are equal
func sortedArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	a_copy := make([]string, len(a))
	b_copy := make([]string, len(b))

	copy(a_copy, a)
	copy(b_copy, b)

	sort.Strings(a_copy)
	sort.Strings(b_copy)

	return reflect.DeepEqual(a_copy, b_copy)
}
