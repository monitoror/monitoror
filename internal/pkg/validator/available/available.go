package available

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/structs"

	"github.com/monitoror/monitoror/api/config/versions"
	pkgValidator "github.com/monitoror/monitoror/internal/pkg/validator"
)

// -------------------------------------------------------
// This file is an helper to validate availability of fields in a given ConfigVersion
// Its purpose is to validate monitorable API params using tags
// Like:
//	type Test struct {
//    Field  string `available:"since:2.2"`
//    Number int `available:"until:3.0"`
//  }
// -------------------------------------------------------

const (
	availableTag    = "available"
	subTagSeparator = ","
)

var (
	sinceSubTagRegexp *regexp.Regexp
	untilSubTagRegexp *regexp.Regexp
)

func init() {
	sinceSubTagRegexp = regexp.MustCompile(`^since=([0-9]+\.[0-9]+)$`)
	untilSubTagRegexp = regexp.MustCompile(`^until=([0-9]+\.[0-9]+)$`)
}

func Struct(s interface{}, version *versions.ConfigVersion) []pkgValidator.Error {
	var errors []pkgValidator.Error

	// Iterate over struct fields but don't go deep inside sub field.
	for _, field := range structs.Fields(s) {
		// Lookup for available tag
		if tagValue := field.Tag(availableTag); tagValue != "" {
			for _, subTag := range strings.Split(tagValue, subTagSeparator) {
				// Lookup for since tag
				if sinceSubTagRegexp.MatchString(subTag) {
					tagVersion := sinceSubTagRegexp.FindStringSubmatch(subTag)[1]
					if version.IsLessThan(versions.RawVersion(tagVersion)) {
						errors = append(errors, &availableError{pkgValidator.ErrorSince, field.Name(), tagVersion})
					}
					continue
				}

				// Lookup for until tag
				if untilSubTagRegexp.MatchString(subTag) {
					tagVersion := untilSubTagRegexp.FindStringSubmatch(subTag)[1]
					if version.IsGreaterThan(versions.RawVersion(tagVersion)) {
						errors = append(errors, &availableError{pkgValidator.ErrorUntil, field.Name(), tagVersion})
					}
					continue
				}

				// Unknown subtag or unsupported version inside validate
				panic(fmt.Sprintf("unknown subtag or unsupported version inside validate. %s", subTag))
			}
		}
	}

	return errors
}
