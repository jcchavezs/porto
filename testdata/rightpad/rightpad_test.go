package rightpad

import "testing"

func TestRightPad(t *testing.T) {
	if want, have := "text ", RightPad("text", 5); want != have {
		t.Errorf("unexpected output, want %q, have %q", want, have)
	}
}
