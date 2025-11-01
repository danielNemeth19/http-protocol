package headers

import (
	"bytes"
	"fmt"
)

var endLine = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	header := make(Headers)
	return header
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if bytes.HasPrefix(data, endLine) {
		return len(endLine), true, nil
	}
	i := bytes.Index(data, endLine)
	if i == -1 {
		return 0, false, nil
	}

	fieldLine := data[:i]
	fieldSep := bytes.IndexByte(fieldLine, ':')
	if fieldSep == -1 {
		return 0, false, fmt.Errorf("Field line supposed to have a ':' separator")
	}

	fieldName := bytes.TrimLeft(fieldLine[:fieldSep], " ")
	fieldValue := bytes.Trim(fieldLine[fieldSep+1:], " ")
	if bytes.HasSuffix(fieldName, []byte(" ")) {
		return 0, false, fmt.Errorf("Invalid field name")
	}
	h[string(fieldName)] = string(fieldValue)
	return len(fieldLine) + len(endLine), false, nil
}
