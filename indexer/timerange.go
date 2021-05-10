package indexer

import (
	"fmt"
	"time"
)

type TimePrecision int

const (
	Year TimePrecision = 1 + iota
	Month
	Day
	Hour
	Minute
	year_layout   = "2006"
	month_layout  = "2006-01"
	day_layout    = "2006-01-02"
	hour_layout   = "2006-01-02 15"
	minute_layout = "2006-01-02 15:04"
)

type TimeRange struct {
	Date      time.Time
	Precision TimePrecision
}

func (r TimeRange) String() string {
	switch r.Precision {
	case Year:
		return r.Date.Format(year_layout)
	case Month:
		return r.Date.Format(month_layout)
	case Day:
		return r.Date.Format(day_layout)
	case Hour:
		return r.Date.Format(hour_layout)
	default:
		return r.Date.Format(minute_layout)
	}
}

func ParseTimeRange(value string) (TimeRange, error) {
	patterns := []struct {
		layout    string
		precision TimePrecision
	}{
		{year_layout, Year},
		{month_layout, Month},
		{day_layout, Day},
		{hour_layout, Hour},
		{minute_layout, Minute},
	}

	for _, pattern := range patterns {
		if len(value) == len(pattern.layout) {
			if date, err := time.Parse(pattern.layout, value); err != nil {
				return TimeRange{}, fmt.Errorf("ParseTimeRange: %w.", err)
			} else {
				return TimeRange{date, pattern.precision}, nil
			}
		}
	}

	return TimeRange{}, fmt.Errorf("ParseTimeRange: Uknown timerange format %q.", value)
}
