package en

import (
	"github.com/raspi/kallu/locales/base"
	"time"
)

// Check implementation
var _ base.Locale = English{}

type English struct {
}

func (f English) GetWeek() string {
	return `week`
}

func (f English) GetWeekDays() map[time.Weekday]string {
	return map[time.Weekday]string{
		time.Monday:    `mon`,
		time.Tuesday:   `tue`,
		time.Wednesday: `wed`,
		time.Thursday:  `thu`,
		time.Friday:    `fri`,
		time.Saturday:  `sat`,
		time.Sunday:    `sun`,
	}
}

func (f English) GetMonths() map[time.Month]string {
	return map[time.Month]string{
		time.January:   `january`, // 1
		time.February:  `february`,
		time.March:     `march`,
		time.April:     `april`,
		time.May:       `may`,
		time.June:      `june`, // 6
		time.July:      `july`,
		time.August:    `august`,
		time.September: `september`,
		time.October:   `october`,
		time.November:  `november`,
		time.December:  `december`, // 12
	}
}
