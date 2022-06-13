package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	var decoded strings.Builder
	var last string
	var escape bool

	for _, x := range encoded {
		cur := string(x)

		if unicode.IsDigit(x) {
			if escape {
				last = cur
				escape = false
				continue
			}

			if last == "" {
				return "", ErrInvalidString
			}

			dig, err := strconv.Atoi(string(x))
			if err != nil {
				return "", ErrInvalidString
			}

			decoded.WriteString(strings.Repeat(last, dig))

			last = ""
			continue
		}

		if last != "" && !escape {
			decoded.WriteString(last)
		}

		if cur == `\` {
			if !escape {
				escape = true
				continue
			}
			escape = false
		}

		last = cur
	}

	if last != "" && !escape {
		decoded.WriteString(last)
	}

	return decoded.String(), nil
}
