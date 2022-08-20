package month

import (
	"fmt"
	"github.com/raspi/kallu/locales/base"
	"strings"
	"time"
	"unicode/utf8"
)

type Month struct {
	m     time.Time
	dow   time.Weekday
	now   time.Time
	start time.Time // Start week
	end   time.Time // End week
}

func New(year int, month time.Month, dow time.Weekday, now time.Time) (m Month) {

	m = Month{
		dow: dow,
		m:   time.Date(year, month, 1, 0, 0, 0, 0, time.Local),
		now: now,
	}

	m.start = m.getStart()
	m.end = m.getEnd()

	return m
}

func (mon Month) GetMonth() (start time.Time) {
	return mon.m
}

func (mon Month) getStart() (start time.Time) {
	start = mon.m

	for start.Weekday() != mon.dow {
		start = start.AddDate(0, 0, -1)
	}

	return start
}

func (mon Month) getEnd() (end time.Time) {
	end = mon.m.AddDate(0, 1, -1)

	if end.Weekday() == mon.dow {
		end = end.AddDate(0, 0, 7)
	} else {
		for end.Weekday() != mon.dow {
			end = end.AddDate(0, 0, 1)
		}
	}

	end = end.AddDate(0, 0, -1)

	return end
}

func (mon Month) GetStartEndWeek() (start time.Time, end time.Time) {
	return mon.start, mon.end
}

func (mon Month) GetDaysWeeks(start time.Time, end time.Time) (weeks int, days int) {
	c := start

	for !c.Equal(end) {
		c = c.AddDate(0, 0, 1)
		days++
	}

	return days / 7, days
}

func (mon Month) GetLastMonth() Month {
	last := mon.m.AddDate(0, -1, 0)
	return New(last.Year(), last.Month(), mon.dow, mon.now)
}

func (mon Month) GetNextMonth() Month {
	last := mon.m.AddDate(0, 1, 0)
	return New(last.Year(), last.Month(), mon.dow, mon.now)
}

func (mon Month) PrintMonth(months []Month, locale base.Locale) {
	const (
		esc             = "\033["
		Clear           = esc + "0m"
		SetForeground   = esc + "38;5;"
		SetBackground   = esc + "48;5;"
		SetUnderlineOn  = esc + "4m"
		SetUnderlineOff = esc + "24m"
		DefaultFG       = SetForeground + "250m"
	)

	weekLocalized := locale.GetWeek()
	weekdaysLocalized := locale.GetWeekDays()
	monthsLocalized := locale.GetMonths()

	monthCount := len(months)

	var headerCol []string
	headerSize := 0

	curr := mon.getStart()

	// "Week" localized
	headerSize += 5
	headerCol = append(headerCol, weekLocalized)

	// Day names
	for di := 0; di < 7; di++ {
		headerSize += 4
		headerCol = append(headerCol, weekdaysLocalized[curr.Weekday()])
		curr = curr.AddDate(0, 0, 1)
	}
	headerSize += 1

	fmt.Print(SetBackground + "238m")
	fmt.Print(SetForeground + "245m")

	// Print month and year header
	for mIdx, m := range months {
		header := fmt.Sprintf(`%v %4v`, monthsLocalized[m.GetMonth().Month()], m.GetMonth().Year())
		if mon.now.Year() == m.m.Year() && mon.now.Month() == m.m.Month() {
			header = "[" + header + "]"
		}

		o := centerString(header, headerSize)

		if mIdx < monthCount-1 {
			o += " | "
		}

		fmt.Print(o)
	}

	fmt.Println(Clear)

	// Day name header
	fmt.Print(SetBackground + "238m")
	fmt.Print(SetForeground + "245m")
	fmt.Print(SetUnderlineOn)

	for mIdx, _ := range months {
		if mIdx > 0 {
			// Separator
			fmt.Print(" | ")
		}

		for idx, hdr := range headerCol {
			switch idx {
			case 0: // Week
				fmt.Printf(`%4s `, hdr)
			default:
				fmt.Printf(`%4s`, hdr)
			}

			if idx == 7 {
				fmt.Print(` `)
			}
		}
	}

	fmt.Println(Clear)

	// Calculate week row count
	maxweeks := 0
	for _, m := range months {
		start, end := m.GetStartEndWeek()
		weeks, _ := m.GetDaysWeeks(start, end)
		if weeks > maxweeks {
			maxweeks = weeks
		}
	}

	for weekIndex := 0; weekIndex < maxweeks+1; weekIndex++ {

		if (weekIndex & 1) == 0 {
			fmt.Print(SetBackground + "235m")
		} else {
			fmt.Print(SetBackground + "236m")
		}

		// Print week and days
		for monthIdx, m := range months {

			_, currweeknum := m.now.ISOWeek()

			start, end := m.GetStartEndWeek()
			if weekIndex > 0 {
				// the magick
				start = start.AddDate(0, 0, 7*weekIndex)
			}

			if start.Before(end) {
				_, weeknum := start.ISOWeek()

				fmt.Print(SetForeground + "245m")

				// Week number
				if m.GetMonth().Month() == start.Month() && m.GetMonth().Year() == m.now.Year() && currweeknum == weeknum {
					// Current week
					fmt.Print(SetForeground + "255m")
				}

				// Print week number
				fmt.Printf(` #%-2v  `, weeknum)

				fmt.Print(DefaultFG)

				prevornext := false

				for i := 0; i < 7; i++ {
					if m.GetMonth().Month() == start.Month() && start.Equal(m.now) {
						fmt.Print(SetForeground + "255m")
						fmt.Print("[")
						fmt.Print(SetUnderlineOn)
					} else {
						fmt.Print(" ")
					}

					if start.Month() != m.GetMonth().Month() && !prevornext {
						// Previous or next month
						fmt.Print(SetForeground + "240m")
						prevornext = true
					}

					if start.Day() == 1 && m.m.Month() == start.Month() {
						fmt.Print(DefaultFG)
						prevornext = false
					}

					// Day number
					fmt.Printf(`%2v`, start.Day())

					if m.GetMonth().Month() == start.Month() && start.Equal(m.now) {
						fmt.Print(SetUnderlineOff)
						fmt.Print("]")
						fmt.Print(DefaultFG)
					} else {
						fmt.Print(" ")
					}

					start = start.AddDate(0, 0, 1)
				}
			} else {
				// Add padding for missing week

				// Week number
				fmt.Print(`      `)
				for i := 0; i < 7; i++ {
					fmt.Print(`    `)
				}

			}

			if monthIdx < monthCount-1 {
				fmt.Print(DefaultFG)
				fmt.Print(" | ")
			}
		}

		fmt.Println(Clear)
	}
}

func centerString(str string, width int) string {
	spaces := int(float64(width-utf8.RuneCountInString(str)) / 2)
	return strings.Repeat(" ", spaces) + str + strings.Repeat(" ", width-(spaces+utf8.RuneCountInString(str)))
}
