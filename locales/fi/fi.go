package fi

import (
	"github.com/raspi/kallu/locales/base"
	"time"
)

// Check implementation
var _ base.Locale = Finnish{}

type Finnish struct {
}

func (f Finnish) GetWeek() string {
	return `vko`
}

func (f Finnish) GetWeekDays() map[time.Weekday]string {
	return map[time.Weekday]string{
		time.Monday:    `ma`,
		time.Tuesday:   `ti`,
		time.Wednesday: `ke`,
		time.Thursday:  `to`,
		time.Friday:    `pe`,
		time.Saturday:  `la`,
		time.Sunday:    `su`,
	}
}

func (f Finnish) GetMonths() map[time.Month]string {
	return map[time.Month]string{
		time.January:   `tammikuu`,
		time.February:  `helmikuu`,
		time.March:     `maaliskuu`,
		time.April:     `huhtikuu`,
		time.May:       `toukokuu`,
		time.June:      `kesäkuu`,
		time.July:      `heinäkuu`,
		time.August:    `elokuu`,
		time.September: `syyskuu`,
		time.October:   `lokakuu`,
		time.November:  `marraskuu`,
		time.December:  `joulukuu`,
	}
}
