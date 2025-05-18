package utils

import (
	"fmt"
	"strings"
)

func Capital(text string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(text[0:1]), text[1:])
}

func GenerateNameWithSpaceName(spaceName string, name string, delimeter string) string {
	if spaceName == "" {
		return name
	}
	return fmt.Sprintf("%s%s%s", spaceName, delimeter, name)
}
