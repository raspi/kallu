package tcell

import "fmt"

const (
	esc             = "\033["
	Clear           = esc + "0m"
	setForeground   = esc + "38;5;"
	setBackground   = esc + "48;5;"
	setUnderlineOn  = esc + "4m"
	setUnderlineOff = esc + "24m"
)

type Option func(a *ANSI)

type ANSI struct {
	BGColor   *uint8
	Underline *bool
	FGColor   *uint8
}

func BG(color uint8) Option {
	return func(a *ANSI) {
		a.BGColor = &color
	}
}

func FG(color uint8) Option {
	return func(a *ANSI) {
		a.FGColor = &color
	}
}

func Ansi(opts ...Option) (a *ANSI) {
	a = &ANSI{}

	if opts != nil {
		for _, apply := range opts {
			// Apply option
			apply(a)
		}
	}

	return a
}

func (a ANSI) String() (s string) {

	if a.FGColor != nil {
		s += fmt.Sprintf(setForeground+`%dm`, *a.FGColor)
	}

	if a.BGColor != nil {
		s += fmt.Sprintf(setBackground+`%dm`, *a.BGColor)
	}

	if a.Underline != nil {
		if *a.Underline {
			s += setUnderlineOn
		} else {
			s += setUnderlineOff
		}
	}

	return s
}

func Underline(enable bool) Option {
	return func(a *ANSI) {
		a.Underline = &enable
	}
}
