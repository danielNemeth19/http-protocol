package headers

import (
	"fmt"
	"strings"

	"github.com/danielNemeth19/http-protocol/internal/request"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if strings.HasPrefix(string(data), request.EndLine) {
		return len(request.EndLine), true, nil
	}
	fieldLine := strings.Split(string(data), request.EndLine)
	if len(fieldLine) < 2 {
		return 0, false, nil
	}
	parts := strings.Split(fieldLine[0], ":")
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("Field line supposed to have two parts, got: %d\n", len(parts))
	}
	fieldName, fieldValue := parts[0], parts[1]
	fmt.Println(fieldName, fieldValue)
	return len(fieldLine), false, nil
}
