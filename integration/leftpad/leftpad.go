// a coment!

package leftpad

// LeftPad pads a string from the left
func LeftPad(s string, length int) string {
	sLength := len(s)
	if sLength >= length {
		return s
	}

	for i := 0; i < sLength-length; i++ {
		s = " " + s
	}

	return s
}
