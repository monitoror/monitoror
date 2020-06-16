package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var removeNull *regexp.Regexp

func init() {
	removeNull = regexp.MustCompile(`\"[^"]+\"\s*:\s*null,?`)
}

func Keys(m interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()
	strKeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strKeys[i] = fmt.Sprintf(`%v`, keys[i])
	}

	return strings.Join(strKeys, ", ")
}

func Stringify(v interface{}) string {
	bytes, _ := json.Marshal(v)
	s := string(bytes)
	s = removeNull.ReplaceAllString(s, ``)
	s = strings.ReplaceAll(s, `,}`, `}`)

	return s
}
