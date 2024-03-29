package month

import (
	"fmt"
	"github.com/raspi/kallu/month/internal/table"
	"github.com/raspi/kallu/month/internal/tcell"
	"github.com/raspi/kallu/month/internal/trow"
	"golang.org/x/text/message"
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

func (mon Month) GetDaysWeeks(start time.Time, end time.Time) (weeks uint, days uint) {
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

func (mon Month) PrintMonth(months []Month, tr *message.Printer, useColor bool) {
	weekdaysLocalized := map[time.Weekday]string{
		time.Monday:    tr.Sprintf(`short.monday`), // mon
		time.Tuesday:   tr.Sprintf(`short.tuesday`),
		time.Wednesday: tr.Sprintf(`short.wednesday`),
		time.Thursday:  tr.Sprintf(`short.thursday`),
		time.Friday:    tr.Sprintf(`short.friday`),
		time.Saturday:  tr.Sprintf(`short.saturday`),
		time.Sunday:    tr.Sprintf(`short.sunday`), // sun
	}

	monthsLocalized := map[time.Month]string{
		time.January:   tr.Sprintf(`january`), // 1
		time.February:  tr.Sprintf(`february`),
		time.March:     tr.Sprintf(`march`),
		time.April:     tr.Sprintf(`april`),
		time.May:       tr.Sprintf(`may`),
		time.June:      tr.Sprintf(`june`), // 6
		time.July:      tr.Sprintf(`july`),
		time.August:    tr.Sprintf(`august`),
		time.September: tr.Sprintf(`september`),
		time.October:   tr.Sprintf(`october`),
		time.November:  tr.Sprintf(`november`),
		time.December:  tr.Sprintf(`december`), // 12
	}

	weekLocalized := tr.Sprintf(`week`)

	DefaultFG := uint8(250)

	monthCount := len(months)

	separator := tcell.New(tcell.Ansi(tcell.FG(245)), `|`)
	emptyC := tcell.New(` `)
	hdrRow := trow.New()
	mtable := table.New(useColor, nil)

	// Week
	weekRowCells := []*tcell.Cell{
		tcell.New(fmt.Sprintf(`%3s `, weekLocalized)),
	}

	// Day names
	curr := mon.getStart()
	for di := 0; di < 7; di++ {
		weekRowCells = append(weekRowCells)
		wdn := tcell.New(fmt.Sprintf(`%3s`, weekdaysLocalized[curr.Weekday()]))

		if di < 6 {
			// Add padding
			wdn.Add(` `)
		}

		weekRowCells = append(weekRowCells, wdn)

		curr = curr.AddDate(0, 0, 1)
	}

	// generate week & day header for each month
	for i := 0; i < monthCount; i++ {
		wrc := weekRowCells

		if i == 0 {
			// Set only on first month
			wrc[0] = tcell.New(
				tcell.Ansi(
					tcell.Underline(true),
					tcell.BG(238),
					tcell.FG(245),
				),
				wrc[0].GetValue()[0],
			)
		}

		hdrRow.Add(wrc...)

		if i < monthCount-1 {
			// Add separator
			hdrRow.Add(separator)
		}

	}
	weekWidth := hdrRow.GetWidths()[0]

	// Add header (week mon tue wed thu fri sat sun separator * monthCount)
	mtable.AddRow(hdrRow)

	// Calculate week row count
	maxweeks := uint(0)
	for _, m := range months {
		start, end := m.GetStartEndWeek()
		weeks, _ := m.GetDaysWeeks(start, end)
		if weeks > maxweeks {
			maxweeks = weeks
		}
	}

	// generate main week view
	for weekIndex := uint(0); weekIndex < maxweeks+1; weekIndex++ {
		wdRow := trow.New()

		// Add weeks and days
		for monthIdx, m := range months {
			_, currweeknum := m.now.ISOWeek()

			start, end := m.GetStartEndWeek()
			if weekIndex > 0 {
				// the magick
				start = start.AddDate(0, 0, int(7*weekIndex))
			}

			if start.Before(end) {
				_, weeknum := start.ISOWeek()

				// Week number cell
				weekC := tcell.New(tcell.Ansi(tcell.FG(245)))

				// Alternate background colors
				if (weekIndex & 1) == 0 {
					weekC.Add(tcell.Ansi(tcell.BG(235)))
				} else {
					weekC.Add(tcell.Ansi(tcell.BG(236)))
				}

				if m.GetMonth().Month() == start.Month() && m.GetMonth().Year() == m.now.Year() && currweeknum == weeknum {
					weekC.Add(tcell.Ansi(tcell.FG(255)))
				}

				// Week number
				weekC.Add(fmt.Sprintf(`%*s`, weekWidth-2, fmt.Sprintf(`#%-2d`, weeknum)))
				weekC.Add(tcell.Ansi(tcell.FG(DefaultFG)))

				// Add week to row
				wdRow.Add(weekC)

				// Previous or next month day
				prevornext := false

				for i := 0; i < 7; i++ {
					// Day number
					dayN := tcell.New()

					currentDay := false

					if m.GetMonth().Month() == start.Month() && start.Equal(m.now) {
						currentDay = true
					} else {
						currentDay = false
					}

					if start.Month() != m.GetMonth().Month() && !prevornext {
						dayN.Add(tcell.Ansi(tcell.FG(240)))

						// Previous or next month
						prevornext = true
					}

					if start.Day() == 1 && m.m.Month() == start.Month() {
						// Start of month
						dayN.Add(tcell.Ansi(tcell.FG(DefaultFG)))

						prevornext = false
					}

					dateFmt := ` %2d `
					if currentDay {
						dateFmt = `[%d]`

						dayN.Add(tcell.Ansi(tcell.FG(255), tcell.Underline(true)))

					}

					dayN.Add(fmt.Sprintf(dateFmt, start.Day()))

					if currentDay {
						dayN.Add(tcell.Ansi(tcell.Underline(false), tcell.FG(DefaultFG)))

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

	// generate top month row
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

		mnameCell := tcell.New(
			tcell.Ansi(
				tcell.BG(238),
				tcell.FG(245),
			),
		)

		spaces := padding / 2
		padding -= spaces

		// Add padding
		mnameCell.Add(strings.Repeat(` `, spaces))

		if isCurrentMonth {
			mnameCell.Add(
				tcell.Ansi(
					tcell.Underline(true),
					tcell.FG(DefaultFG),
				),
			)
		}

		// Add header
		mnameCell.Add(header)

		if isCurrentMonth {
			mnameCell.Add(
				tcell.Ansi(
					tcell.Underline(false),
					tcell.FG(245),
				),
			)
		}

		// Add padding
		mnameCell.Add(strings.Repeat(` `, padding))

		// Month name
		monthRow.Add(mnameCell)

		// Add separator
		if i < monthCount-1 {
			monthRow.Add(separator)
		}

	}

	// Month name table
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
