package tracker

import (
	"strconv"
	"time"
)

// AddZeroTime adds "0" in front of single digit val, for example "7" becomes "07". Assumes 0 <= val <= 60
func AddZeroTime(val int) string {
	if val < 10 {
		return "0" + strconv.Itoa(val)
	}
	return strconv.Itoa(val)
}

// TextifyTime returns current time as: "{Weekday} / {Year} {Month} {Day} / {Hour}:{Minute}:{Second}"
func TextifyTime() string {
	year, month, day := time.Now().Date()
	weekday := time.Now().Weekday()
	hour, min, sec := time.Now().Clock()
	hourText := AddZeroTime(hour)
	minText := AddZeroTime(min)
	secText := AddZeroTime(sec)
	return weekday.String() + " / " + strconv.Itoa(year) + " " + month.String() + " " + strconv.Itoa(day) + " / " + hourText + ":" + minText + ":" + secText
}
