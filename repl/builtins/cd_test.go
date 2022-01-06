package builtins

import "testing"

func TestHandleChange_defaultCase(t *testing.T) {
	got, _ := handleChange("/test/1", "2")
	if got != "/test/1/2" {
		t.Fail()
	}
}

func TestHandleChange_upADirectory(t *testing.T) {
	got, _ := handleChange("/test/1", "..")
	if got != "/test" {
		t.Fail()
	}
}

func TestHandleChange_relativeDirectory(t *testing.T) {
	got, _ := handleChange("/test/1", "./2")
	if got != "/test/1/2" {
		t.Fail()
	}
}

func TestHandleChange_up2Directories(t *testing.T) {
	got, _ := handleChange("/test/1/2/3", "../../")
	if got != "/test/1" {
		t.Fail()
	}
}
