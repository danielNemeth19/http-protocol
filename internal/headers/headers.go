package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

var endLine = []byte("\r\n")

var specialChars = []string{"!", "#", "$", "%", "&", "'", "*", "+", "-", ".", "^", "_", "`", "|", "~"}

type Headers map[string]string

func NewHeaders() Headers {
	header := make(Headers)
	return header
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if bytes.HasPrefix(data, endLine) {
		return len(endLine), true, nil
	}
	endLineSep := bytes.Index(data, endLine)
	if endLineSep == -1 {
		return 0, false, nil
	}

	fieldLine := data[:endLineSep]
	fieldSep := bytes.IndexByte(fieldLine, ':')
	if fieldSep == -1 {
		return 0, false, fmt.Errorf("Field line supposed to have a ':' separator")
	}

	fieldName := string(fieldLine[:fieldSep])
	fieldValue := string(fieldLine[fieldSep+1:])

	fieldName = strings.TrimLeft(fieldName, " ")
	fieldValue = strings.Trim(fieldValue, " ")

	if strings.HasSuffix(fieldName, " ") {
		return 0, false, fmt.Errorf("Invalid field name")
	}
	for _, c := range fieldName {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || slices.Contains(specialChars, string(c))) {
			return 0, false, fmt.Errorf("Invalid header key: key %s contains invalid char %s", fieldName, string(c))
		}
	}
	fieldName = strings.ToLower(fieldName)
	h[fieldName] = fieldValue
	return len(fieldLine) + len(endLine), false, nil
}
