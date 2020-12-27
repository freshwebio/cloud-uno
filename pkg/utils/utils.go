package utils

import "strings"

// CommaSeparatedListContains determines whether the provided
// comma-separated list contains the provided search string as an exact match
// of a comma-separated value.
func CommaSeparatedListContains(commaSeparatedList string, search string) bool {
	i := 0
	found := false
	list := strings.Split(commaSeparatedList, ",")
	for !found && i < len(list) {
		found = list[i] == search
		i = i + 1
	}
	return found
}
