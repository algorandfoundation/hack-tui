package style

import "testing"

func Test_WithStyles(t *testing.T) {
	testStr := Red.Render("Amazing") + Green.Render("World")
	render := WithControls("", testStr)
	if render != testStr {
		t.Error("Should be empty")
	}
	render = WithOverlay("", testStr)
	if render != testStr {
		t.Error("Should be empty")
	}

	render = TruncateLeft(testStr, 7)
	if render != "World" {
		t.Error("Should be World")
	}
}