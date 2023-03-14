package tcell

import (
	"fmt"
	"unicode"
)

// Cell is single cell that can be inserted into trow.Row
type Cell struct {
	width uint
	value []any // strings and ANSI colors etc. which creates a view
}

func New(v ...any) (c *Cell, err error) {
	c = &Cell{
		width: 0,
	}

	if v != nil {
		err = c.add(v...)
	}

	return c, err
}

func (c *Cell) add(v ...any) error {
	for _, vt := range v {
		c.value = append(c.value, vt)

		switch vt.(type) {
		case *ANSI:

		case string:
			s := vt.(string)

			for ci, rc := range s {

				switch rc {
				case '\n', '\v', '\f', '\r', 0x85, 0xA0:
					return fmt.Errorf(`string contains new line at position %d`, ci)
				case '\t':
					return fmt.Errorf(`string contains tab at position %d`, ci)
				}

				if !unicode.IsPrint(rc) {
					return fmt.Errorf(`string contains non-printable character %q at position %d`, rc, ci)
				}
			}

			// Update width
			c.width += uint(len(s))
		default:
			// Not supported
			return fmt.Errorf(`not supported: v=%[1]T vt=%[2]T`, v, vt)
		}

	}

	return nil
}

// Add adds new strings and/or ANSI to cell
func (c *Cell) Add(v ...any) (err error) {
	if v != nil {
		err = c.add(v...)
	}

	return err
}

func (c *Cell) GetWidth() uint {
	return c.width
}

func (c *Cell) String() (s string) {
	for _, vt := range c.value {
		switch vt.(type) {
		case string:
			s += vt.(string)
		}
	}

	return s
}

// GetValue returns internal representation of Cell values
func (c *Cell) GetValue() []any {
	return c.value
}
