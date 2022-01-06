package display

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	log "k8s.io/klog"

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

func sortDescending(g *Grid) []*Cell {
	tmp := make([]*Cell, len(g.cells))
	copy(tmp, g.cells)
	sort.Slice(tmp, func(i, j int) bool {
		return g.cells[i].width > g.cells[j].width
	})
	return tmp
}

func (s *Screen) guessNumLines(max int) int {
	guessMinNum := 0
	totalSoFar := 0
	numCells := len(s.grid.cells)
	cells := sortDescending(s.grid)
	for _, c := range cells {
		log.Info("Current total: ", c.width+totalSoFar)
		if c.width+totalSoFar <= max {
			guessMinNum += 1
			totalSoFar += c.width
		} else {
			guessMaxLines := numCells / guessMinNum
			log.Info("Guess: ", guessMaxLines, numCells, guessMinNum)
			if numCells%guessMinNum != 0 {
				guessMaxLines += 1
			}
			return guessMaxLines
		}
		totalSoFar += s.grid.options.padding.Width()
	}
	return 1
}

func (s *Screen) getCellWidths() []int {
	widths := make([]int, len(s.grid.cells))
	for i, c := range s.grid.cells {
		widths[i] = c.width
	}
	return widths
}

func (s *Screen) getWidestCell() []int {
	widest := 0
	for _, c := range s.grid.cells {
		if c.width > widest {
			widest = c.width
		}
	}
	return []int{widest}
}

func sumWidths(widths []int) int {
	sum := 0
	for _, w := range widths {
		sum += w
	}
	return sum
}

func (s *Screen) getColumnWidths(numLines int, numColumns int) *Dimensions {
	// log.Infof("Lines: %d; Columns: %d", numLines, numColumns)
	widths := make([]int, numColumns)
	for i, c := range s.grid.cells {
		var index int
		switch s.grid.options.direction {
		case LeftToRight:
			index = i % numColumns
		case TopToBottom:
			index = (i / numLines) % numColumns
		}
		// log.Infof("i: %d; index: %d", i, index)
		widthAndPadding := c.width + s.grid.options.padding.Width()
		widths[index] = max(widths[index], widthAndPadding)
	}
	return &Dimensions{
		widths: widths,
		lines:  numLines,
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
		return nil, errors.New("cannot fit grid into screen smaller than widest cell")
	}

	numCells := len(s.grid.cells)
	if numCells == 0 {
		return &Dimensions{
			lines:  0,
			widths: make([]int, 0),
		}, nil
	}
	if numCells == 1 {
		return &Dimensions{
			lines:  1,
			widths: []int{s.grid.cells[0].width},
		}, nil
	}
	guessNumLines := s.guessNumLines(max)
	log.V(2).Infoln("Guess num lines", guessNumLines)

	if guessNumLines == 1 {
		return &Dimensions{
			lines:  1,
			widths: s.getCellWidths(),
		}, nil
	}

	smallestDimensions := &Dimensions{
		lines:  len(s.grid.cells),
		widths: s.getWidestCell(),
	}
	for lines := len(s.grid.cells) / 2; lines > 0; lines-- {
		numColumns := numCells / lines
		log.V(2).Infof("numColumns: %d; remainder: %d", numColumns, numCells%lines)

		numLines := lines
		if numCells%lines != 0 {
			numLines++
		}

		totalSeparatorWidth := (numColumns - 1) * s.grid.options.padding.Width()
		log.V(2).Infof("totalSeparatorWidth: %d; max: %d", totalSeparatorWidth, max)

		if max < totalSeparatorWidth {
			continue
		}

		adjustedWidth := max - totalSeparatorWidth
		potentialDimensions := s.getColumnWidths(numLines, numColumns)

		totalSum := sumWidths(potentialDimensions.widths)
		log.V(2).Infof("adjustedWidth: %d; widths: %v; totalSum: %d", adjustedWidth, potentialDimensions.widths, totalSum)

		if totalSum < adjustedWidth {
			smallestDimensions = potentialDimensions
		} else {
			continue
		}
	}
	log.V(2).Infoln("Smalles dimensions", smallestDimensions)
	return smallestDimensions, nil
}

func (s *Screen) Print() error {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	log.V(2).Infoln("Width", width)
	dimensions, err := s.dimensionsFromWidth(width)
	log.V(2).Infoln("dimensions, err", dimensions, err)
	if err != nil {
		return err
	}
	log.V(2).Infoln("Dimension widths", dimensions.widths)
	numWidths := len(dimensions.widths)
	log.V(2).Infoln("num widths", numWidths)
	numCells := len(s.grid.cells)
	log.V(2).Infoln("num cells", numCells)
	log.V(2).Infof("Options: %v", s.grid.options)
	for y := 0; y < dimensions.lines; y++ {
		for x := 0; x < numWidths; x++ {
			index := 0
			switch s.grid.options.direction {
			case LeftToRight:
				index = y*numWidths + x
			case TopToBottom:
				index = y + dimensions.lines*x
			}
			// fmt.Infof("|x: %d, y: %d, index: %d, numCells: %d\n", x, y, index, numCells)
			if index >= numCells {
				// log.Infoln("Adding newline")
				// fmt.Infof("\n")
				break
			}

			cell := s.grid.cells[index]
			if x == numWidths-1 {
				switch cell.alignment {
				case Left:
					fmt.Printf("%s\n", cell.content)
				case Right:
					extraSpaces := dimensions.widths[x] - cell.width
					fmt.Printf("%s%s\n", strings.Repeat(" ", max(extraSpaces, 0)), cell.content)
				}
			} else {
				filling := s.grid.options.padding.Padding(cell.alignment)
				// log.Infof("Filling: '%s'", filling)
				extraSpaces := dimensions.widths[x] - cell.width - len(filling)
				// log.Infof("Extra spaces: %d, widths: %v, width: %d", extraSpaces, dimensions.widths[x], cell.width)
				switch cell.alignment {
				case Left:
					fmt.Printf("%s%s%s", cell.content, strings.Repeat(" ", max(0, extraSpaces)), filling)
				case Right:
					fmt.Printf("%s%s%s", strings.Repeat(" ", max(0, extraSpaces)), cell.content, filling)
				}
			}
		}
	}
	fmt.Println("")
	return nil
}
