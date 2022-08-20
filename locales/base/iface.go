package base

import "time"

type Locale interface {
	GetWeekDays() map[time.Weekday]string // mon-sun
	GetMonths() map[time.Month]string     // jan-dec
	GetWeek() string                      // `week` localized
}
