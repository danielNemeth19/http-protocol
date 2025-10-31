package headers

import (
	"bytes"
	"fmt"
)

var endLine = []byte("\r\n")

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if bytes.HasPrefix(data, endLine) {
		return len(endLine), true, nil
	}
	i := bytes.Index(data, endLine); if i == -1 {
		return 0, false, nil
	}
	fieldLine := data[:i]
	fieldSep := bytes.IndexByte(fieldLine, ':'); if fieldSep == -1 {
		return 0, false, fmt.Errorf("Field line supposed to a ':' separator")
	}


	// parts := strings.Split(fieldLine[0], ":")
	// if len(parts) != 2 {
		// return 0, false, fmt.Errorf("Field line supposed to have two parts, got: %d\n", len(parts))
	// }
	// fieldName, fieldValue := parts[0], parts[1]
	// fmt.Println(fieldName, fieldValue)
	return 0, false, nil
}
