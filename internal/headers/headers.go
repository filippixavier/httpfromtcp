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

func (h Headers) Get(key string) string {
	k := strings.ToLower(key)

	if val, ok := h[k]; ok {
		return val
	}
	return ""
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

	invalidReg, err := regexp.Compile("[^a-zA-Z0-9!#$%&'*+-.^_`|~]")

	if err != nil {
		return 0, false, err
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	if invalidReg.MatchString(key) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}
	h.Set(key, string(value))
	return idx + 2, false, nil
}
