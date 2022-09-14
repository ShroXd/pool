package crawler

import (
	"testing"
)

func TestIncreaseNum(t *testing.T) {
	got, _ := increaseNum("1")
	want := "2"

	if want != got {
		t.Errorf("Expected '%s', but got '%s'", want, got)
	}
}
