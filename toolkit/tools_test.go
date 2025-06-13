package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var testtools Tools

	s := testtools.RandomString(7)
	if len(s) != 7 {
		t.Error("wrong length of random string returned")
	}
}
