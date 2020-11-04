package leftpad

import "testing"

func TestLeftPad(t *testing.T) {
	if want, have := " text", LeftPad("text", 5); want != have {
		t.Errorf("unexpected output, want %q, have %q", want, have)
	}
}
