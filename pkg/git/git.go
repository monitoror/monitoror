package git

import (
	"strings"
)

func HumanizeBranch(branch string) string {
	// Remove refs/head
	branch = strings.Replace(branch, "refs/heads/", "", 1)

	return branch
}
