package utils

import "time"

func datesEqual(t1 time.Time, t2 time.Time) bool {
	var equal bool = true
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	if y1 != y2 {
		equal = false
	}
	if m1 != m2 {
		equal = false
	}
	if d1 != d2 {
		equal = false
	}
	return equal
}

func boolToInt(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
