package display

import (
	"reflect"
	"testing"
)

func defaultGrid() *Grid {
	g := NewGridWithOptions(&GridOptions{padding: NewWhitespacePadding(1)})
	g.AddCell(NewCell("a"))
	g.AddCell(NewCell("b"))
	g.AddCell(NewCell("c"))
	g.AddCell(NewCell("d"))

	return g
}

func TestScreen_guessNumLines_withDefaultGrid(t *testing.T) {
	s := NewScreen(defaultGrid())

	got := s.guessNumLines(1)
	if got != 4 {
		t.Log("Should have guessed 4 lines, got ", got)
		t.Fail()
	}

	got = s.guessNumLines(2)
	if got != 4 {
		t.Log("Should have guessed 4 lines, got ", got)
		t.Fail()
	}

	got = s.guessNumLines(3)
	if got != 2 {
		t.Log("Should have guessed 2 lines, got ", got)
		t.Fail()
	}

	got = s.guessNumLines(6)
	if got != 2 {
		t.Log("Should have guessed 2 lines, got ", got)
		t.Fail()
	}

	got = s.guessNumLines(7)
	if got != 1 {
		t.Log("Should have guessed 1 lines, got ", got)
		t.Fail()
	}
}
func TestScreen_guessNumLines_withVariantGrid(t *testing.T) {
	s := NewScreen(variantGrid())

	got := s.guessNumLines(9)
	if got != 2 {
		t.Log("Should have guessed 2 lines, got ", got)
		t.Fail()
	}
}

func TestScreen_getCellWidths(t *testing.T) {
	s := NewScreen(defaultGrid())

	got := s.getCellWidths()
	want := []int{1, 1, 1, 1}
	if !reflect.DeepEqual(got, want) {
		t.Fail()
	}
}

func variantGrid() *Grid {
	g := NewGrid()
	g.AddCell(NewCell("a"))
	g.AddCell(NewCell("aa"))
	g.AddCell(NewCell("aaa"))
	g.AddCell(NewCell("aaaa"))
	return g
}

func TestScreen_getColumnWidths_LeftToRight(t *testing.T) {
	s := NewScreen(variantGrid())

	got := s.getColumnWidths(4, 1)
	want := &Dimensions{widths: []int{6}, lines: 4}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(2, 3)
	want = &Dimensions{widths: []int{6, 4, 5}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(2, 2)
	want = &Dimensions{widths: []int{5, 6}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(1, 4)
	want = &Dimensions{widths: []int{3, 4, 5, 6}, lines: 1}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(1, 3)
	want = &Dimensions{widths: []int{6, 4, 5}, lines: 1}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	s.grid.options.direction = TopToBottom
	got = s.getColumnWidths(2, 2)
	want = &Dimensions{widths: []int{4, 6}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}
}

func TestScreen_getColumnWidths_TopToBottom(t *testing.T) {
	s := NewScreen(variantGrid())
	s.grid.options.direction = TopToBottom

	got := s.getColumnWidths(4, 1)
	want := &Dimensions{widths: []int{6}, lines: 4}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(2, 3)
	want = &Dimensions{widths: []int{4, 6, 0}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(2, 2)
	want = &Dimensions{widths: []int{4, 6}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(1, 4)
	want = &Dimensions{widths: []int{3, 4, 5, 6}, lines: 1}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(1, 3)
	want = &Dimensions{widths: []int{6, 4, 5}, lines: 1}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}

	got = s.getColumnWidths(2, 2)
	want = &Dimensions{widths: []int{4, 6}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got %v, wanted %v", got, want)
		t.Fail()
	}
}

func TestScreen_dimensionsFromWidth(t *testing.T) {
	s := NewScreen(variantGrid())
	got, err := s.dimensionsFromWidth(1)
	if err == nil || got != nil {
		t.Log("Error not returned when max is less than widest cell")
		t.Fail()
	}

	s2 := NewScreen(NewGrid())
	got, err = s2.dimensionsFromWidth(1)
	want := &Dimensions{widths: []int{}, lines: 0}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Log("Did not handle empty case")
		t.Fail()
	}

	g3 := NewGrid()
	g3.AddCell(NewCell("test"))
	s3 := NewScreen(g3)
	got, err = s3.dimensionsFromWidth(5)
	want = &Dimensions{widths: []int{4}, lines: 1}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Log("Did not handle case of only 1 cell")
		t.Fail()
	}

	got, err = s.dimensionsFromWidth(24)
	want = &Dimensions{widths: []int{1, 2, 3, 4}, lines: 1}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Log("Did not handle case of guessed 1 line")
		t.Fail()
	}

	got, _ = s.dimensionsFromWidth(14)
	want = &Dimensions{widths: []int{5, 6}, lines: 2}
	if !reflect.DeepEqual(got, want) {
		t.Logf("Got: %v, wanted: %v", got, want)
		t.Fail()
	}
}
