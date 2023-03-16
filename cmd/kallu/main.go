package main

import (
	"flag"
	"fmt"
	"github.com/raspi/kallu/month"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"strings"
	"time"
)

//go:generate gotext -srclang=en update -out generated_locales.go -lang=en,fi-FI

// These are set with Makefile -X=main.VERSION, etc
var (
	VERSION   = `v0.0.0`
	BUILD     = `dev`
	BUILDDATE = `0000-00-00T00:00:00+00:00`
)

const (
	AUTHOR   = `Pekka Järvinen`
	HOMEPAGE = `https://github.com/raspi/kallu`
	YEAR     = 2022
)

/*
getEnv fetches environmental variables
- NO_COLOR disable color output?
- LANG which language to use
*/
func getEnv() (useColor bool, locale language.Tag) {
	useColor = true

	// Default locale
	locale = language.English

	for _, k := range os.Environ() {
		arr := strings.SplitAfter(k, `=`)
		key, value := strings.TrimRight(arr[0], `=`), arr[1]

		switch key {
		case `LANG`:
			if value == `C` {
				continue
			}

			// fi_FI => fi-FI
			value = strings.Replace(value, `_`, `-`, -1)
			value = strings.ToLower(value)
			if strings.Contains(value, `.`) {
				// remove .UTF-8
				value = value[:strings.LastIndex(value, `.`)]
			}

			candidates, _, err := language.ParseAcceptLanguage(value)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, `error: could not parse environmental parameter LANG value %q`+"\n", value)
				break // Fail and use default language
			}

			if len(candidates) > 0 {
				// Set language to first one
				locale = candidates[0]
			}

		case `NO_COLOR`:
			// See: https://no-color.org/
			useColor = false
		}

	}

	return
}

func main() {
	tmpnow := time.Now()
	now := time.Date(tmpnow.Year(), tmpnow.Month(), tmpnow.Day(), 0, 0, 0, 0, tmpnow.Location())

	envColor, locale := getEnv()
	tr := message.NewPrinter(locale)

	useColor := true

	if !envColor {
		useColor = false
	}

	howManyNext := flag.Uint(`next`, 1, tr.Sprintf(`How many next months`))
	howManyPrev := flag.Uint(`prev`, 1, tr.Sprintf(`How many previous months`))

	howManyChunks := flag.Uint(`count`, 3, tr.Sprintf(`How many months per line`))

	selectedYear := flag.Uint(`year`, uint(now.Year()), tr.Sprintf(`Year (defaults to current year)`))
	selectedMonth := flag.Uint(`month`, uint(now.Month()), tr.Sprintf(`Month 1-12 (defaults to current month)`))

	fullYear := flag.Bool(`fullyear`, false, tr.Sprintf(`Print full year`))
	currentOnly := flag.Bool(`one`, false, tr.Sprintf(`Only one month, equivalent to -next 0 -prev 0`))

	noColor := flag.Bool(`no-color`, false, tr.Sprintf(`Disable color output`))

	selectedDow := flag.Uint(`dow`, uint(time.Monday), tr.Sprintf(`Start day for week 0-6 (sun-sat)`))

	flag.Usage = func() {
		f := os.Args[0] // exe name
		_, _ = fmt.Fprintf(os.Stdout, `kallu - simple CLI calendar`+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `Version %v %v %v`+"\n", VERSION, BUILD, BUILDDATE)
		_, _ = fmt.Fprintf(os.Stdout, `(c) %v %v- [ %v ]`+"\n", AUTHOR, YEAR, HOMEPAGE)
		_, _ = fmt.Fprintf(os.Stdout, "\n")

		_, _ = fmt.Fprintf(os.Stdout, tr.Sprintf(`Parameters:`)+"\n")

		paramMaxLen := 0

		flag.VisitAll(func(f *flag.Flag) {
			l := len(f.Name)
			if l > paramMaxLen {
				paramMaxLen = l
			}
		})

		flag.VisitAll(func(f *flag.Flag) {
			padding := strings.Repeat(` `, paramMaxLen-len(f.Name))

			displayDef := true
			defValue := f.DefValue

			switch f.Name {
			case `month`, `year`:
				displayDef = false
				defValue = ``
			}

			_, _ = fmt.Fprintf(os.Stdout, `  -%s%s   %s   `, f.Name, padding, f.Usage)

			if displayDef {
				_, _ = fmt.Fprintf(os.Stdout, tr.Sprintf(`default: %q`, defValue))
			}

			_, _ = fmt.Fprintf(os.Stdout, "\n")

		})

		_, _ = fmt.Fprintf(os.Stdout, "\n")

		_, _ = fmt.Fprintf(os.Stdout, tr.Sprintf(`Examples:`)+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `  `+tr.Sprintf(`Full year:`)+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `    %s -fullyear`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, `  `+tr.Sprintf(`Only this month:`)+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `    %s -one`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, `  `+tr.Sprintf(`One calendar at a time:`)+"\n")
		_, _ = fmt.Fprintf(os.Stdout, `    %s -count 1`+"\n", f)
		_, _ = fmt.Fprintf(os.Stdout, "\n")
	}

	flag.Parse()

	if *noColor {
		useColor = false
	}

	if *selectedMonth > 12 || *selectedMonth == 0 {
		_, _ = fmt.Fprintf(os.Stderr, tr.Sprintf(`invalid month: %d`, *selectedMonth))
		os.Exit(1)
	}

	if *selectedDow > 6 {
		_, _ = fmt.Fprintf(os.Stderr, tr.Sprintf(`invalid starting day of week: %d`, *selectedDow))
		os.Exit(1)
	}

	dow := time.Weekday(*selectedDow)

	currentMonth := month.New(int(*selectedYear), time.Month(*selectedMonth), dow, now)

	if *howManyChunks == 0 {
		_, _ = fmt.Fprintf(os.Stderr, tr.Sprintf(`count must be > 0`))
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

	chunks := *howManyChunks
	currChunk := uint(0)

	var monthList []month.Month
	for _, m := range months {
		monthList = append(monthList, m)
		currChunk++

		if currChunk == chunks {
			currentMonth.PrintMonth(monthList, tr, useColor)
			fmt.Println()

			// Clear list for next chunk
			monthList = []month.Month{}
			currChunk = 0
		}
	}

	if len(monthList) > 0 {
		// Print remaining month(s)
		currentMonth.PrintMonth(monthList, tr, useColor)
		fmt.Println()
	}

}
