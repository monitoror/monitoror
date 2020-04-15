package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

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
	return string(bytes)
}
