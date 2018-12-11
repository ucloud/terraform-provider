package ucloud

import (
	"bytes"
	"fmt"
	"strings"
)

// Converter is use for converting string to another string with specifical style
type Converter interface {
	convert(string) (string, error)
	unconvert(string) (string, error)
	mustConvert(string) string
	mustUnconvert(string) string
}

type styleConverter struct {
	input    map[string]string
	reversed map[string]string
}

func newStyleConverter(input map[string]string) *styleConverter {
	return &styleConverter{
		input:    input,
		reversed: buildReversedStringMap(input),
	}
}

type upperConverter struct {
	*styleConverter
}

func newUpperConverter(specials map[string]string) *upperConverter {
	return &upperConverter{
		styleConverter: newStyleConverter(
			specials,
		),
	}
}

// convert is an utils used for converting upper case name with underscore into lower case with underscore.
func (cvt *upperConverter) convert(input string) (string, error) {
	if input != strings.ToUpper(input) {
		return "", fmt.Errorf("excepted input string is uppercase with underscore, got %s", input)
	}
	return cvt.mustConvert(input), nil
}

func (cvt *upperConverter) mustConvert(input string) string {
	return strings.ToLower(input)
}

// unconvert is an utils used for converting lower case with underscore into upper case name with underscore.
func (cvt *upperConverter) unconvert(input string) (string, error) {
	if input != strings.ToLower(input) {
		return "", fmt.Errorf("excepted input string is lowercase with underscore, got %s", input)
	}
	return strings.ToUpper(input), nil
}

func (cvt *upperConverter) mustUnconvert(input string) string {
	return strings.ToUpper(input)
}

type lowerCamelConverter struct {
	*styleConverter
}

func newLowerCamelConverter(specials map[string]string) *lowerCamelConverter {
	return &lowerCamelConverter{
		styleConverter: newStyleConverter(
			specials,
		),
	}
}

func (cvt *lowerCamelConverter) convert(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	if 'A' <= input[0] && input[0] <= 'Z' {
		return "", fmt.Errorf("excepted lower camel should not be leading by uppercase character, got %s", input)
	}

	return lowerCamelToLower(input), nil
}

func (cvt *lowerCamelConverter) mustConvert(input string) string {
	output, _ := cvt.convert(input)
	return output
}

func (cvt *lowerCamelConverter) unconvert(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	if input != strings.ToLower(input) {
		return "", fmt.Errorf("excepted input string is lowercase with underscore, got %s", input)
	}

	return cvt.mustUnconvert(input), nil
}

func (cvt *lowerCamelConverter) mustUnconvert(input string) string {
	return lowerToLowerCamel(input)
}

type upperCamelConverter struct {
	*styleConverter
}

func newUpperCamelConverter(specials map[string]string) *upperCamelConverter {
	return &upperCamelConverter{
		styleConverter: newStyleConverter(
			specials,
		),
	}
}

func (cvt *upperCamelConverter) convert(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	if 'a' <= input[0] && input[0] <= 'z' {
		return "", fmt.Errorf("excepted upper camel should not be leading by lowercase character, got %s", input)
	}

	return lowerCamelToLower(strings.ToLower(input[:1]) + input[1:]), nil
}

func (cvt *upperCamelConverter) mustConvert(input string) string {
	output, _ := cvt.convert(input)
	return output
}

func (cvt *upperCamelConverter) unconvert(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	if input != strings.ToLower(input) {
		return "", fmt.Errorf("excepted input string is lowercase with underscore, got %s", input)
	}

	output := lowerToLowerCamel(input)
	return strings.ToUpper(output[:1]) + output[1:], nil
}

func (cvt *upperCamelConverter) mustUnconvert(input string) string {
	output, _ := cvt.unconvert(input)
	return output
}

func lowerCamelToLower(input string) string {
	// eg. createFail -> create_fail; createUDBFAIL -> create_udb_fail -> createUdbFail
	var state int
	var words []string
	buf := strings.Builder{}
	for i := 0; i < len(input); i++ {
		c, l1 := input[i], lookAhead(&input, i, 1)

		// last character
		if l1 == 0 {
			buf.Write(bytes.ToLower([]byte{c}))
			words = append(words, buf.String())
			buf.Reset()
			break
		}

		if state == 0 {
			if 'A' <= l1 && l1 <= 'Z' {
				// createing UDBInstance
				//         ^ ^
				//         | |
				//         c l1
				buf.WriteByte(c)
				state = 1

				words = append(words, buf.String())
				buf.Reset()
			} else {
				// createi ngUDBInstance
				//       ^ ^
				//       | |
				//       c l1
				buf.WriteByte(c)
			}

			continue
		}

		if state == 1 {
			if 'A' <= l1 && l1 <= 'Z' {
				// createingU DBInstance
				//          ^ ^
				//          | |
				//          c l1
				buf.WriteByte(c + ('a' - 'A'))
				state = 3
			} else {
				// createingI nstance
				//          ^ ^
				//          | |
				//          c l1
				buf.WriteByte(c + ('a' - 'A'))
				state = 0
			}

			continue
		}

		if state == 3 {
			if 'A' <= l1 && l1 <= 'Z' {
				// createingUD BInstance
				//           ^ ^
				//           | |
				//           c l1
				buf.WriteByte(c + ('a' - 'A'))
			} else {
				// createingUDBI nstance
				//             ^ ^
				//             | |
				//             c l1
				words = append(words, buf.String())
				buf.Reset()

				buf.WriteByte(c + ('a' - 'A'))
				state = 0
			}

			continue
		}
	}

	return strings.Join(words, "_")
}

func lowerToLowerCamel(input string) string {
	iL := strings.Split(input, "_")
	oL := make([]string, len(iL))
	for i, s := range iL {
		oL[i] = strings.Title(s)
	}
	output := strings.Join(oL, "")
	return strings.ToLower(output[:1]) + output[1:]
}

func lookAhead(input *string, index, forward int) byte {
	if len((*input)) <= index+forward {
		return 0
	}
	return (*input)[index+forward]
}
