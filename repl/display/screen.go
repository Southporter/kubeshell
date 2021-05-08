package display

import (
  "errors"
  "os"
  
  "golang.org/x/term"
)

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
  guessMinNum := 0
  totalSoFar := 0
  numCells := len(s.grid.cells)
  cells := reverse(s.grid.cells)
  for _, c := range cells {
    if c.width + totalSoFar <= max {
      guessMinNum += 1;
      totalSoFar += c.width
    } else {
      guessMaxLines := numCells / guessMinNum
      if numCells % guessMinNum != 0 {
        guessMaxLines += 1
      }
      return guessMaxLines
    }
    totalSoFar += s.grid.options.padding.Width()
  }
  return 1
}

func reverse(cells []*Cell) []*Cell {
  numCells := len(cells)
  reversed := make([]*Cell, numCells)
  for i, c := range cells {
    reversed[numCells - i - 1] = c
  }
  return reversed
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
  return sum
}

func (s *Screen) getColumnWidths(numLines int, numColumns int) *Dimensions {
  widths := make([]int, numColumns)
  for i, c := range s.grid.cells {
    index := i
    switch s.grid.options.direction {
    case LeftToRight:
      index = i % numColumns
    case TopToBottom:
      index = i / numLines
    }
    widths[index] = max(widths[index], c.width)
  }
  return &Dimensions{
    widths: widths,
    lines: numLines,
  }
}

func max(a int, b int) int {
  if a > b {
    return a
  }
  return b
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
    potentialDimensions := s.getColumnWidths(lines, numColumns)

    if sumWidths(potentialDimensions.widths) < adjustedWidth {
      smallestDimensions = potentialDimensions
    } else {
      break
    }
  }
  return smallestDimensions, nil
}

func (s *Screen) Print() error {
  width, _, err := term.GetSize(int(os.Stdin.Fd()))
  if err != nil {
    return err
  }
  dimensions, err := s.dimensionsFromWidth(width)
  if err != nil {
    return err
  }
  numWidths := len(dimensions.widths)
  numCells := len(s.grid.cells)
  for y := 0; y < dimensions.lines; y++ {
    for x := 0; x < numWidths; x++ {
      index := 0
      switch s.grid.options.direction {
      case LeftToRight:
        index = y * numWidths + x
      case TopToBottom:
        index = y + dimensions.lines * x
      }
      if index >= numCells {
        continue
      }

      cell := s.grid.cells[index]
      if x == numWidths - 1 {

      }
    }
  }
  return nil
}
