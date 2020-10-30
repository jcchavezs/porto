package rightpad // import "wrong/import/rightpad"

// RightPad pads a string from the right
func RightPad(s string, length int) string {
	sLength := len(s)
	if sLength >= length {
		return s
	}

	for i := 0; i < length-sLength; i++ {
		s += " "
	}

	return s
}
