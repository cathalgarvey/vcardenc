package vcardenc

// finds the end of a quoted string, assuming the opening quotation mark is
// stripped from line.
func parseQuotedValue(line string, delimCs []rune, expectClosing bool) (parsedLine, remaining string, err error) {
	var (
		escaped     bool
		parsedChars []rune
	)
	for n, c := range line {
		if !escaped && c == '\\' {
			escaped = true
			continue
		}
		if runeSliceContains(delimCs, c) {
			if escaped {
				if c == 'n' {
					c = '\n'
				}
				parsedChars = append(parsedChars, c)
				escaped = false
				continue
			}
			parsedLine = line[:n]
			if expectClosing {
				remaining = line[n+1:]
			} else {
				remaining = line[n:]
			}
			return string(parsedChars), remaining, nil
		}
		parsedChars = append(parsedChars, c)
		escaped = false
	}
	if !expectClosing {
		return string(parsedChars), remaining, nil
	}
	return "", "", ErrFailedToParseQuotedString
}

func runeSliceContains(slice []rune, R rune) bool {
	for _, r := range slice {
		if r == R {
			return true
		}
	}
	return false
}

func stringSliceContains(slice []string, S string) bool {
	for _, s := range slice {
		if s == S {
			return true
		}
	}
	return false
}
