// Adapted from https://raw.githubusercontent.com/dustin/go-humanize/master/times.go

package format

import (
	"fmt"
	"sort"
	"time"
)

// Seconds-based time units
const (
	Day      = 24 * time.Hour
	Week     = 7 * Day
	Month    = 30 * Day
	Year     = 12 * Month
	LongTime = 37 * Year
)

// Time formats a time into a relative string.
//
// Time(someT) -> "3 weeks ago"
func Time(then time.Time) string {
	return relTime(then.UTC(), time.Now().UTC(), "ago", "from now")
}

// A relTimeMagnitude struct contains a relative time point at which
// the relative format of time will switch to a new format string.  A
// slice of these in ascending order by their "D" field is passed to
// CustomRelTime to format durations.
//
// The Format field is a string that may contain a "%s" which will be
// replaced with the appropriate signed label (e.g. "ago" or "from
// now") and a "%d" that will be replaced by the quantity.
//
// The DivBy field is the amount of time the time difference must be
// divided by in order to display correctly.
//
// e.g. if D is 2*time.Minute and you want to display "%d minutes %s"
// DivBy should be time.Minute so whatever the duration is will be
// expressed in minutes.
type relTimeMagnitude struct {
	D     time.Duration
	Fn    func(diff time.Duration, lbl string) string
	DivBy time.Duration
}

var defaultMagnitudes = []relTimeMagnitude{
	{time.Second, func(diff time.Duration, lbl string) string { return "just now" }, time.Second},
	{2 * time.Second, func(diff time.Duration, lbl string) string { return fmt.Sprintf("1s %s", lbl) }, 1},
	{time.Minute, func(diff time.Duration, lbl string) string { return fmt.Sprintf("%ds %s", diff, lbl) }, time.Second},
	{2 * time.Minute, func(diff time.Duration, lbl string) string { return fmt.Sprintf("1m %s", lbl) }, 1},
	{time.Hour, func(diff time.Duration, lbl string) string { return fmt.Sprintf("%dm %s", diff, lbl) }, time.Minute},
	{2 * time.Hour, func(diff time.Duration, lbl string) string { return fmt.Sprintf("1h %s", lbl) }, 1},
	{Day, func(diff time.Duration, lbl string) string { return fmt.Sprintf("%dh %s", diff, lbl) }, time.Hour},
	{2 * Day, func(diff time.Duration, lbl string) string { return fmt.Sprintf("1d %s", lbl) }, 1},
	{Week, func(diff time.Duration, lbl string) string { return fmt.Sprintf("%dd %s", diff, lbl) }, Day},
	{2 * Week, func(diff time.Duration, lbl string) string { return fmt.Sprintf("1w %s", lbl) }, 1},
	{Month, func(diff time.Duration, lbl string) string { return fmt.Sprintf("%dw %s", diff, lbl) }, Week},
}

// RelTime(timeInPast, timeInFuture, "earlier", "later") -> "3 weeks earlier"
func relTime(a, b time.Time, albl, blbl string) string {
	return customRelTime(a, b, albl, blbl, defaultMagnitudes)
}

func customRelTime(a, b time.Time, albl, blbl string, magnitudes []relTimeMagnitude) string {
	lbl := albl
	diff := b.Sub(a)

	if a.After(b) {
		lbl = blbl
		diff = a.Sub(b)
	}

	// Find the largest magnitude
	largestMagnitude := magnitudes[len(magnitudes)-1].D

	// If the difference is greater than the largest magnitude, format the date in local time
	if diff >= largestMagnitude {
		return a.Local().Format("Jan 2 2006")
	}

	n := sort.Search(len(magnitudes), func(i int) bool {
		return magnitudes[i].D > diff
	})

	// If no magnitude is large enough, use the largest magnitude available
	if n >= len(magnitudes) {
		n = len(magnitudes) - 1
	}
	mag := magnitudes[n]

	if mag.DivBy == 1 {
		return mag.Fn(diff, lbl)
	}

	return mag.Fn(diff/mag.DivBy, lbl)
}
