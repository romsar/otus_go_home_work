package hw02unpackstring

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	var decoded strings.Builder
	var last string
	var escaping bool

	for _, x := range encoded {
		cur := string(x)

		if unicode.IsDigit(x) {
			if escaping {
				last = cur
				escaping = false
				continue
			}

			if last == "" {
				return "", ErrInvalidString
			}

			dig, err := strconv.Atoi(cur)
			if err != nil {
				return "", errors.Wrap(ErrInvalidString, err.Error())
			}

			decoded.WriteString(strings.Repeat(last, dig))

			last = ""
			escaping = false

			continue
		}

		if last != "" && !escaping {
			decoded.WriteString(last)
		}

		isBackSlash := cur == `\`

		if escaping && !isBackSlash {
			last = `\` + cur
			escaping = false
			continue
		}

		escaping = !escaping && isBackSlash
		last = cur
	}

	if last != "" && !escaping {
		decoded.WriteString(last)
	}

	return decoded.String(), nil
}
