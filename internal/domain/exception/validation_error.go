package exception

import (
	"bytes"
	"strings"
)

type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < len(ve); i++ {
		buff.WriteString(ve[i].Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}
