package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {

	var buf string
	var esc bool
	res := strings.Builder{}
	for _, b := range str {
		letter := fmt.Sprintf("%c", b)
		digit, err := strconv.Atoi(letter)
		if letter == `\` && !esc {
			esc = true
		} else if err != nil && letter != `\` && esc {
			return "", ErrInvalidString
		} else if err != nil || esc {
			res.Write([]byte(buf))
			buf = letter
			if esc {
				esc = false
			}
		} else {
			if buf == "" {
				return "", ErrInvalidString
			} else {
				res.Write([]byte(strings.Repeat(buf, digit)))
				buf = ""
			}
		}
	}
	if esc {
		return "", ErrInvalidString
	}
	res.Write([]byte(buf))
	return res.String(), nil
}
