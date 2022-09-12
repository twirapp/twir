package helpers

import (
	"fmt"
	"strings"
	"time"
)

func dateDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func Duration(t time.Time) string {
	date := strings.Builder{}
	y, m, d, h, mi, s := dateDiff(t, time.Now())
	if y > 0 {
		fmt.Fprintf(&date, "%dy ", y)
	}

	if m > 0 {
		fmt.Fprintf(&date, "%dmo ", m)
	}

	if d > 0 {
		fmt.Fprintf(&date, "%dd ", d)
	}

	if h > 0 {
		fmt.Fprintf(&date, "%dh ", h)
	}

	if mi > 0 {
		fmt.Fprintf(&date, "%dm ", mi)
	}

	if s > 0 {
		fmt.Fprintf(&date, "%ds", s)
	}

	return date.String()
}
