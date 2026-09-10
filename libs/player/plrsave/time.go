package plrsave

import "time"

func AsTime(t time.Time) Time {
	return Time{
		Year:         uint16(t.Year()),
		Month:        uint16(t.Month()),
		DayOfWeek:    uint16(t.Weekday()),
		Day:          uint16(t.Day()),
		Hour:         uint16(t.Hour()),
		Minute:       uint16(t.Minute()),
		Second:       uint16(t.Second()),
		Milliseconds: uint16(t.Nanosecond() / 1e6),
	}
}

type Time struct {
	Year         uint16
	Month        uint16
	DayOfWeek    uint16
	Day          uint16
	Hour         uint16
	Minute       uint16
	Second       uint16
	Milliseconds uint16
}

func (ts *Time) String() string {
	return ts.Time().String()
}

func (ts *Time) GoString() string {
	return ts.Time().GoString()
}

func (ts *Time) Time() time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Date(
		int(ts.Year), time.Month(ts.Month), int(ts.Day),
		int(ts.Hour), int(ts.Minute), int(ts.Second),
		int(ts.Milliseconds)*int(time.Millisecond),
		time.Local,
	)
}
