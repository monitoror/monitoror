package git

import (
	"fmt"
	"strings"
)

func HumanizeBranch(branch string) string {
	// Remove refs/head
	branch = strings.Replace(branch, "refs/heads/", "", 1)

	// Add #
	if !strings.HasPrefix(branch, "#") {
		branch = fmt.Sprintf("#%s", branch)
	}
	return branch
}
