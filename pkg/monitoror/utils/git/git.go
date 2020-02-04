package git

import (
	"fmt"
	"strings"
)

const branchPrefix = "@"

func HumanizeBranch(branch string) string {
	// Remove refs/head
	branch = strings.Replace(branch, "refs/heads/", "", 1)

	// Add @
	if !strings.HasPrefix(branch, branchPrefix) {
		branch = fmt.Sprintf("%s%s", branchPrefix, branch)
	}

	return branch
}
