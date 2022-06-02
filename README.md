# kallu

![GitHub All Releases](https://img.shields.io/github/downloads/raspi/kallu/total?style=for-the-badge)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/raspi/kallu?style=for-the-badge)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/raspi/kallu?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/raspi/kallu)](https://goreportcard.com/report/github.com/raspi/kallu)


Simple CLI calendar.

![Screenshot](https://github.com/raspi/kallu/blob/main/_assets/kallu_default.png)

![Screenshot](https://github.com/raspi/kallu/blob/main/_assets/kallu_count_next.png)

## Usage

```text
kallu - simple CLI calendar
Version v0.0.0 dev 0000-00-00T00:00:00+00:00
(c) Pekka Järvinen 2022- [ https://github.com/raspi/kallu ]

Parameters:
  -count      How many months per line   default: "3"
  -dow        Start day for week 0-6 (sun-sat)   default: "1"
  -fullyear   Print full year   default: "false"
  -month      Month 1-12 (defaults to current month)   default: "6"
  -next       How many next months   default: "1"
  -one        Only one month, equivalent to -next 0 -prev 0   default: "false"
  -prev       How many previous months   default: "1"
  -year       Year (defaults to current year)   default: "2022"

Examples:
  Full year:
    ./kallu -fullyear
  Only this month:
    ./kallu -one
    - equivalent to ./kallu -next 0 -prev 0
  One calendar at a time:
    ./kallu -count 1
```

## Is it any good?

Yes.
