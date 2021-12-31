package display

import (
	"testing"
)

func TestNewCell(t *testing.T) {
	got := NewCell("test")

	if got.width != 4 {
		t.Log("Width was not 4")
		t.Fail()
	}

	got = NewCell("pneumonoultramicroscopicsilicovolcanoconiosis")

	if got.width != 45 {
		t.Log("Width was not 45")
		t.Fail()
	}
}

func TestPaddingWidth(t *testing.T) {
	got := NewTextPadding("--")

	if got.Width() != 2 {
		t.Log("Width incorrectly calculated. Expected 2, got ", got.Width())
		t.Fail()
	}

	got.spaces = 1
	if got.Width() != 3 {
		t.Log("Width does not add spaces correctly. Expected 3, got ", got.Width())
		t.Fail()
	}

	got = NewWhitespacePadding(4)
	if got.Width() != 4 {
		t.Log("Whitespace padding not calculated correctly. Expected 4, got ", got.Width())
		t.Fail()
	}
}

func TestPaddingWithAlignment(t *testing.T) {
	padding := NewWhitespacePadding(2)

	got := padding.Padding(Left)
	if got != "  " {
		t.Fail()
	}

	padding.text = "*"

	got = padding.Padding(Left)
	if got != "  *" {
		t.Fail()
	}

	got = padding.Padding(Right)
	if got != "*  " {
		t.Fail()
	}

	padding = NewTextPadding("*")
	got = padding.Padding(Right)
	if got != "*" {
		t.Fail()
	}
}

func TestDimensionsTotalWidth(t *testing.T) {
	d := &Dimensions{lines: 5, widths: []int{1, 2, 1, 2, 1, 2}}
	p := NewWhitespacePadding(1)

	got := d.TotalWidth(p)
	if got != 9+5 {
		t.Fail()
	}

	p = NewTextPadding("--")

	got = d.TotalWidth(p)
	if got != 9+10 {
		t.Fail()
	}
}
