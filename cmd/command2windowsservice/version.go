package main

import (
	"fmt"
	"strings"
)

var (
	version string
	commit  string
	branch  string
)

func getVersion() string {
	parts := []string{"server"}
	if version != "" {
		parts = append(parts, version)
	} else {
		parts = append(parts, "unknown")
	}
	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}
	return strings.Join(parts, " ")
}
