package display

import "errors"

type Screen struct {
  grid *Grid
}

func NewScreen(g *Grid) *Screen {
  return &Screen{
    grid: g,
  }
}

func (s *Screen) FitToWidth(max int) error {
  return nil
}

func (s *Screen) guessNumLines(max int) int {
  return max
}

func getCellWidths(cells []*Cell) []int {
  widths := make([]int, len(cells))
  for i, c := range(cells) {
    widths[i] = c.width
  }
  return widths
}

func sumWidths(widths []int) int {
  sum := 0
  for _, w := range widths {
    sum += w;
  }
  return w
}

func getColumnWidths(numLines int, numColumns int) *Dimensions {
  return &Dimensions{
    widths: []int{},
    lines: numLines,
  }
}

func (s *Screen) dimensionsFromWidth(max int) (*Dimensions, error) {
  if s.grid.widestCell > max {
    return nil, errors.New("Cannot fit grid into screen smaller than widest cell")
  }


  numCells := len(s.grid.cells)
  if numCells == 0 {
    return &Dimensions{
      lines: 0,
      widths: make([]int, 0),
    }, nil
  }
  if numCells == 1 {
    return &Dimensions{
      lines: 1,
      widths: []int{s.grid.cells[0].width},
    }, nil
  }
  guessNumLines := s.guessNumLines(max)

  if guessNumLines == 1 {
    return &Dimensions{
      lines: 1,
      widths: getCellWidths(s.grid.cells),
    }, nil
  }

  var smallestDimensions *Dimensions
  for lines := guessNumLines; lines > 0; lines-- {
    numColumns := numCells / lines

    if numCells % lines != 0 {
      numColumns += 1
    }

    totalSeparatorWidth := (numColumns - 1) * s.grid.options.padding.Width()

    if max < totalSeparatorWidth {
      continue
    }

    adjustedWidth := max - totalSeparatorWidth;
    potentialDimensions := getColumnWidths(lines, numColumns)

    if sumWidths(potentialDimensions.widths) < adjustedWidth {
      smallestDimensions = potentialDimensions
    } else {
      break
    }
  }
  return smallestDimensions, nil
}


