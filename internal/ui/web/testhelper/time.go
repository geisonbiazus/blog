package testhelper

import "time"

func ParseTime(timeString string) time.Time {
	t, _ := time.Parse(time.RFC3339, "2021-04-03T00:00:00+00:00")
	return t
}
