package utils

import (
	"bytes"
)

type MultiError []error

func (errs MultiError) Error() string {
	if len(errs) == 0 {
		return ""
	}

	buffer := new(bytes.Buffer)
	for i, err := range errs {
		if i >= 1 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(err.Error())
	}
	return buffer.String()
}
