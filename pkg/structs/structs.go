package structs

import (
	"strings"

	"github.com/fatih/structs"
)

func GetJSONFieldName(field *structs.Field) string {
	return strings.Split(field.Tag("json"), ",")[0]
}

func GetQueryFieldName(field *structs.Field) string {
	return strings.Split(field.Tag("query"), ",")[0]
}
