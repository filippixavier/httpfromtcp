package headers

import (
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	headers := make(map[string]string)

	return headers
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	str := string(data)

	lines := strings.Split(str, "\r\n")
	read := 0

	invalidReg, err := regexp.Compile("[^a-zA-Z0-9!#$%&'*+-.^_`|~]")

	if err != nil {
		return 0, false, err
	}

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, ": ")

		if len(parts) != 2 {
			return read, false, fmt.Errorf("invalid header field")
		}

		if strings.Contains(parts[0], " ") || strings.Contains(parts[1], " ") {
			return read, false, fmt.Errorf("invalid space header")
		}

		if invalidReg.MatchString(parts[0]) {
			return read, false, fmt.Errorf("header contain invalid charaters")
		}

		key := strings.ToLower(parts[0])

		if val, ok := h[key]; ok {
			h[key] = val + ", " + parts[1]
		} else {
			h[key] = parts[1]
		}

		read += len(line) + 2
	}

	return read, false, nil
}
