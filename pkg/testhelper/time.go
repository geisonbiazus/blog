package testhelper

import "time"

func ParseTime(timeString string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeString)
	return t
}
