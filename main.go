package main

import (
	"flag"
	"fmt"
	"github.com/raspi/kallu/month"
	"os"
	"strings"
	"time"
)

var (
	// These are set with Makefile -X=main.VERSION, etc
	VERSION   = `v0.0.0`
	BUILD     = `dev`
	BUILDDATE = `0000-00-00T00:00:00+00:00`
)

const (
	AUTHOR   = `Pekka Järvinen`
	HOMEPAGE = `https://github.com/raspi/kallu`
	YEAR     = 2022
)

func main() {

	tmpnow := time.Now()
	now := time.Date(tmpnow.Year(), tmpnow.Month(), tmpnow.Day(), 0, 0, 0, 0, tmpnow.Location())

	howManyNext := flag.Uint(`next`, 1, `How many next months`)
	howManyPrev := flag.Uint(`prev`, 1, `How many previous months`)

	howManyChunks := flag.Uint(`count`, 3, `How many months per line`)

	selectedYear := flag.Uint(`year`, uint(now.Year()), `Year (defaults to current year)`)
	selectedMonth := flag.Uint(`month`, uint(now.Month()), `Month 1-12 (defaults to current month)`)

	fullYear := flag.Bool(`fullyear`, false, `Print full year`)
	currentOnly := flag.Bool(`one`, false, `Only one month, equivalent to -next 0 -prev 0`)

	selectedDow := flag.Uint(`dow`, uint(time.Monday), `Start day for week 0-6 (sun-sat)`)

	flag.Usage = func() {
		f := os.Args[0]
		_, _ = fmt.Fprintf(os.Stdout, `kallu - simple CLI calendar`+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `Version %v %v %v`+"\n", VERSION, BUILD, BUILDDATE)
		_, _ = fmt.Fprintf(os.Stdout, `(c) %v %v- [ %v ]`+"\n", AUTHOR, YEAR, HOMEPAGE)
		_, _ = fmt.Fprintf(os.Stdout, "\n")

		_, _ = fmt.Fprintf(os.Stdout, "Parameters:\n")

		paramMaxLen := 0

		flag.VisitAll(func(f *flag.Flag) {
			l := len(f.Name)
			if l > paramMaxLen {
				paramMaxLen = l
			}
		})

		flag.VisitAll(func(f *flag.Flag) {
			padding := strings.Repeat(` `, paramMaxLen-len(f.Name))
			_, _ = fmt.Fprintf(os.Stdout, "  -%s%s   %s   default: %q\n", f.Name, padding, f.Usage, f.DefValue)
		})

		_, _ = fmt.Fprintf(os.Stdout, "\n")

		_, _ = fmt.Fprintf(os.Stdout, `Examples:`+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `  Full year:`+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `    %s -fullyear`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, `  Only this month:`+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `    %s -one`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, `    - equivalent to %s -next 0 -prev 0`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, `  One calendar at a time:`+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `    %s -count 1`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, "\n")
	}

	flag.Parse()

	if *selectedMonth > 12 || *selectedMonth == 0 {
		_, _ = fmt.Fprintf(os.Stderr, `invalid month: %d`, *selectedMonth)
		os.Exit(1)
	}

	if *selectedDow > 6 {
		_, _ = fmt.Fprintf(os.Stderr, `invalid starting day of week: %d`, *selectedDow)
		os.Exit(1)
	}

	dow := time.Weekday(*selectedDow)

	currentMonth := month.New(int(*selectedYear), time.Month(*selectedMonth), dow, now)

	if *howManyChunks == 0 {
		_, _ = fmt.Fprintf(os.Stderr, `count must be > 0`)
		os.Exit(1)
	}

	var months []month.Month

	if *fullYear {
		// Get entire year
		y := month.New(int(*selectedYear), time.January, dow, now)

		for i := 0; i < 12; i++ {
			months = append(months, y)
			y = y.GetNextMonth()
		}

	} else {
		if *currentOnly {
			*howManyPrev = 0
			*howManyNext = 0
		}

		if *howManyPrev > 0 {
			// Add N previous month(s)
			var tmpmonths []month.Month

			lastMonth := currentMonth
			for i := uint(0); i < *howManyPrev; i++ {
				lastMonth = lastMonth.GetLastMonth()
				tmpmonths = append(tmpmonths, lastMonth)
			}

			for i := len(tmpmonths) - 1; i > -1; i-- {
				// Add in reverse
				months = append(months, tmpmonths[i])
			}

		}

		months = append(months, currentMonth)

		if *howManyNext > 0 {
			// Add N next month(s)
			nextMonth := currentMonth
			for i := uint(0); i < *howManyNext; i++ {
				nextMonth = nextMonth.GetNextMonth()
				months = append(months, nextMonth)
			}
		}
	}

	chunks := int(*howManyChunks)
	currChunk := 0

	var monthList []month.Month
	for _, m := range months {
		monthList = append(monthList, m)
		currChunk++

		if currChunk == chunks {
			currentMonth.PrintMonth(monthList)
			fmt.Println()
			monthList = []month.Month{}
			currChunk = 0
		}
	}

	if len(monthList) > 0 {
		// Print remaining month(s)
		currentMonth.PrintMonth(monthList)
		fmt.Println()
	}

}
