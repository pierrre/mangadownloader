package mangadownloader

import (
	"bytes"
)

type MultiError []error

func (errs MultiError) Error() string {
	buffer := new(bytes.Buffer)
	for i, err := range errs {
		if i >= 1 {
			buffer.WriteString("\n")
		}
		buffer.WriteString(err.Error())
	}
	return buffer.String()
}
