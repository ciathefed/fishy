package preprocessor

import "unicode"

func skipWhitespace(s string, i *int) {
	for *i < len(s) && s[*i] == ' ' {
		(*i)++
	}
}

func skipNonWhitespace(s string, i *int) {
	for *i < len(s) && s[*i] != ' ' {
		(*i)++
	}
}

func skipToBreak(s string, i *int) {
	breakChars := []rune{',', ' '}
	for *i < len(s) {
		ch := rune(s[*i])
		if contains(breakChars, ch) {
			break
		}
		(*i)++
	}
}

func skipAlpha(s string, i *int) {
	for *i < len(s) && unicode.IsLetter(rune(s[*i])) {
		(*i)++
	}
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || (ch >= '0' && ch <= '9')
}

func contains(slice []rune, item rune) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func isValidLabelName(label string) bool {
	if len(label) == 0 {
		return false
	}

	firstChar := rune(label[0])
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		return false
	}

	for _, ch := range label[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}

	return true
}
