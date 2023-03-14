package trow

import (
	"github.com/raspi/kallu/month/internal/tcell"
)

// Row consists of multiple tcell.Cell items
type Row struct {
	columns      []*tcell.Cell
	columnCount  uint
	columnWidths []uint
}

// New creates new row with N columns from tcell.Cell items
func New(row ...*tcell.Cell) (r Row) {
	r = Row{}

	if row != nil {
		r.add(row...)
	}

	return r
}

func (r *Row) Add(row ...*tcell.Cell) {
	if row != nil {
		r.add(row...)
	}
}

func (r *Row) add(row ...*tcell.Cell) {
	for _, cell := range row {
		r.columnCount++
		r.columnWidths = append(r.columnWidths, cell.GetWidth())
		r.columns = append(r.columns, cell)
	}
}

func (r *Row) GetWidths() []uint {
	return r.columnWidths
}

func (r *Row) GetCells() (s []*tcell.Cell) {
	return r.columns
}
