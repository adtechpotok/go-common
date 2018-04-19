package help

import "testing"

func TestLowerCaseFirstLetter(t *testing.T) {
	if m := LowerCaseFirstLetter("Test"); m != "test" {
		t.Errorf("Letter was not lowered, res %s",m)
	}

	if LowerCaseFirstLetter("")!= "" {
		t.Errorf("Empty string error")
	}
}
