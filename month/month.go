package month

import (
	"fmt"
	"github.com/raspi/kallu/month/internal/table"
	"github.com/raspi/kallu/month/internal/tcell"
	"github.com/raspi/kallu/month/internal/trow"
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
	next := mon.m.AddDate(0, 1, 0)
	return New(next.Year(), next.Month(), mon.dow, mon.now)
}

func (mon Month) PrintMonth(months []Month, weekdaysLocalized map[time.Weekday]string, weekLocalized string, monthsLocalized map[time.Month]string, useColor bool) {
	DefaultFG := uint8(250)

	monthCount := len(months)

	separator, err := tcell.New(tcell.Ansi(tcell.FG(245)), `|`)
	if err != nil {
		panic(err)
	}

	emptyC, err := tcell.New(` `)
	if err != nil {
		panic(err)
	}

	hdrRow := trow.New()

	mtable := table.New(useColor, nil)

	for i := 0; i < monthCount; i++ {

		weekName, err := tcell.New(
			tcell.Ansi(
				tcell.Underline(true),
				tcell.BG(238),
				tcell.FG(245),
			),
			weekLocalized+`  `,
		)
		if err != nil {
			panic(err)
		}

		hdrRow.Add(weekName)

		curr := mon.getStart()

		// Day names
		for di := 0; di < 7; di++ {
			wdn, err := tcell.New(fmt.Sprintf(`%3s`, weekdaysLocalized[curr.Weekday()]))
			if err != nil {
				panic(err)
			}

			if di < 6 {
				// Add padding
				err = wdn.Add(` `)
				if err != nil {
					panic(err)
				}
			}

			hdrRow.Add(wdn)

			curr = curr.AddDate(0, 0, 1)
		}

		if i < monthCount-1 {
			// Add separator
			hdrRow.Add(separator)
		}

	}
	weekWidth := hdrRow.GetWidths()[0]

	// Add header (week mon tue wed thu fri sat sun separator * monthCount)
	mtable.AddRow(hdrRow)

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
		wdRow := trow.New()

		// Add weeks and days
		for monthIdx, m := range months {
			_, currweeknum := m.now.ISOWeek()

			start, end := m.GetStartEndWeek()
			if weekIndex > 0 {
				// the magick
				start = start.AddDate(0, 0, 7*weekIndex)
			}

			if start.Before(end) {
				_, weeknum := start.ISOWeek()

				// Week number
				weekC, err := tcell.New(tcell.Ansi(tcell.FG(245)))
				if err != nil {
					panic(err)
				}

				if (weekIndex & 1) == 0 {
					err = weekC.Add(tcell.Ansi(tcell.BG(235)))
					if err != nil {
						panic(err)
					}
				} else {
					err = weekC.Add(tcell.Ansi(tcell.BG(236)))
					if err != nil {
						panic(err)
					}
				}

				if m.GetMonth().Month() == start.Month() && m.GetMonth().Year() == m.now.Year() && currweeknum == weeknum {
					err = weekC.Add(tcell.Ansi(tcell.FG(255)))
					if err != nil {
						panic(err)
					}
				}

				// Week number
				err = weekC.Add(fmt.Sprintf(`%*s`, weekWidth-2, fmt.Sprintf(`#%-2d`, weeknum)))
				if err != nil {
					panic(err)
				}
				wdRow.Add(weekC)

				err = weekC.Add(tcell.Ansi(tcell.FG(DefaultFG)))
				if err != nil {
					panic(err)
				}

				// Previous or next month day
				prevornext := false

				for i := 0; i < 7; i++ {

					// Day number
					dayN, err := tcell.New()
					if err != nil {
						panic(err)
					}

					currentDay := false

					if m.GetMonth().Month() == start.Month() && start.Equal(m.now) {
						currentDay = true
					} else {
						currentDay = false
					}

					if start.Month() != m.GetMonth().Month() && !prevornext {

						err = dayN.Add(tcell.Ansi(tcell.FG(240)))
						if err != nil {
							panic(err)
						}

						// Previous or next month
						prevornext = true
					}

					if start.Day() == 1 && m.m.Month() == start.Month() {
						// Start of month
						err = dayN.Add(tcell.Ansi(tcell.FG(DefaultFG)))
						if err != nil {
							panic(err)
						}

						prevornext = false
					}

					dateFmt := ` %2d `
					if currentDay {
						dateFmt = `[%d]`

						err = dayN.Add(tcell.Ansi(tcell.FG(255), tcell.Underline(true)))
						if err != nil {
							panic(err)
						}

					}

					err = dayN.Add(fmt.Sprintf(dateFmt, start.Day()))
					if err != nil {
						panic(err)
					}

					if currentDay {
						err = dayN.Add(tcell.Ansi(tcell.Underline(false), tcell.FG(DefaultFG)))
						if err != nil {
							panic(err)
						}

						currentDay = false
					}

					start = start.AddDate(0, 0, 1)

					wdRow.Add(dayN)
				}
			} else {
				// Add padding for missing week

				wdRow.Add(emptyC) // week
				wdRow.Add(emptyC) // mon
				wdRow.Add(emptyC) // tue
				wdRow.Add(emptyC) // wed
				wdRow.Add(emptyC) // thu
				wdRow.Add(emptyC) // fri
				wdRow.Add(emptyC) // sat
				wdRow.Add(emptyC) // sun

			}

			// Add separator
			if monthIdx < monthCount-1 {
				wdRow.Add(separator)
			}
		}

		mtable.AddRow(wdRow)
	}

	requiredWidth := 0
	// Get width from title row
	for _, w := range mtable.GetWidth()[0:8] {
		requiredWidth += int(w)
	}

	monthRow := trow.New()
	for i, m := range months {

		isCurrentMonth := false
		header := fmt.Sprintf(
			`%v %4v`,
			monthsLocalized[m.GetMonth().Month()], m.GetMonth().Year(),
		)
		if mon.now.Year() == m.m.Year() && mon.now.Month() == m.m.Month() {
			isCurrentMonth = true
			header = "[" + header + "]"
		}

		padding := requiredWidth - utf8.RuneCountInString(header)

		mnameCell, err := tcell.New(
			tcell.Ansi(
				tcell.BG(238),
				tcell.FG(245),
			),
		)
		if err != nil {
			panic(err)
		}

		spaces := padding / 2
		padding -= spaces

		err = mnameCell.Add(strings.Repeat(` `, spaces))
		if err != nil {
			panic(err)
		}

		if isCurrentMonth {
			err = mnameCell.Add(
				tcell.Ansi(
					tcell.Underline(true),
					tcell.FG(DefaultFG)),
			)
			if err != nil {
				panic(err)
			}

		}

		// Add header
		err = mnameCell.Add(header)
		if err != nil {
			panic(err)
		}

		if isCurrentMonth {
			err = mnameCell.Add(
				tcell.Ansi(
					tcell.Underline(false),
					tcell.FG(245)),
			)
			if err != nil {
				panic(err)
			}
		}

		err = mnameCell.Add(strings.Repeat(` `, padding))
		if err != nil {
			panic(err)
		}

		monthRow.Add(mnameCell)

		// Add separator
		if i < monthCount-1 {
			monthRow.Add(separator)
		}

	}

	mnametable := table.New(useColor, &monthRow)

	// Print month names
	for _, x := range mnametable.GetRows()[0] {
		fmt.Printf(`%v`, x)
	}

	fmt.Println(tcell.Clear)

	// Print weeks
	for _, c := range mtable.GetRows() {
		for _, x := range c {
			fmt.Printf(`%v`, x)
		}

		fmt.Println(tcell.Clear)
	}

}
