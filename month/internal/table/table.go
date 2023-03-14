package table

import (
	"fmt"
	"github.com/raspi/kallu/month/internal/tcell"
	"github.com/raspi/kallu/month/internal/trow"
	"strings"
)

type ColoredTable struct {
	rows         []trow.Row
	columnWidths []uint
	columnCount  int
	useColor     bool
}

func New(useColor bool, row *trow.Row) (ct *ColoredTable) {
	ct = &ColoredTable{
		useColor: useColor,
	}

	if row != nil {
		ct.addRow(*row)
	}

	return ct
}

func (ct *ColoredTable) AddRow(row trow.Row) {
	ct.addRow(row)
}

func (ct *ColoredTable) addRow(row trow.Row) {
	widths := row.GetWidths()
	lw := len(widths)

	if ct.columnCount == 0 {
		// First row sets initial column sizes
		ct.columnCount = lw
		ct.columnWidths = widths
	} else {
		if ct.columnCount != lw {
			panic(fmt.Errorf(`have %d columns, got %d`, ct.columnCount, lw))
		}
	}

	for i, colWidth := range ct.columnWidths {
		if widths[i] > colWidth {
			// Update width
			ct.columnWidths[i] = widths[i]
		}

	}

	ct.rows = append(ct.rows, row)
}

func (ct *ColoredTable) GetColumnCount() int {
	return ct.columnCount
}

func (ct *ColoredTable) GetRows() (s [][]string) {
	for _, r := range ct.rows {
		var n []string

		for cellIndex, cell := range r.GetCells() {

			txt := ``     // Text representation
			outputs := `` // Whole representation with color(s)

			for _, cellValue := range cell.GetValue() {
				switch cellValue.(type) {
				case string:
					vs := cellValue.(string)
					txt += vs
					outputs += vs
				case *tcell.ANSI:
					if ct.useColor {
						outputs += cellValue.(*tcell.ANSI).String()
					}
				}
			}

			outputs += strings.Repeat(` `, int(ct.columnWidths[cellIndex])-len(txt))
			n = append(n, outputs)
		}

		s = append(s, n)

	}

	return s
}

func (ct *ColoredTable) GetWidth() []uint {
	return ct.columnWidths
}
