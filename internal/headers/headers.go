package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	headers := make(map[string]string)

	return headers
}

func (h Headers) Set(key, value string) {
	k := strings.ToLower(key)

	if val, ok := h[k]; ok {
		h[k] = val + ", " + value
	} else {
		h[k] = value
	}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// the empty line
		// headers are done, consume the CRLF
		return 2, true, nil
	}

	str := string(data)

	lines := strings.Split(str, "\r\n")

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
			return 0, false, fmt.Errorf("invalid header field")
		}

		if strings.Contains(parts[0], " ") || strings.Contains(parts[1], " ") {
			return 0, false, fmt.Errorf("invalid space header")
		}

		if invalidReg.MatchString(parts[0]) {
			return 0, false, fmt.Errorf("header contain invalid charaters")
		}

		h.Set(parts[0], parts[1])
	}

	return idx + 2, false, nil
}
