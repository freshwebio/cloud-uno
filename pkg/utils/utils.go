// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package utils

import (
	"os"
	"strings"
)

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

// IsRunningInDockerContainer determines whether or not the current
// program is running in docker.
func IsRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
