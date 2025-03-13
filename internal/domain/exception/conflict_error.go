package exception

import (
	"bytes"
	"strings"
)

type ConflictErrors struct {
	Errors []FieldError
}

func (ce ConflictErrors) Error() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < len(ce.Errors); i++ {
		buff.WriteString(ce.Errors[i].Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

func (ce *ConflictErrors) Add(field, tag string) {
	ce.Errors = append(ce.Errors, FieldError{Field: field, Tag: tag})
}
