package indexer_test

import (
	"testing"
	"time"

	"github.com/cosaques/algolia/indexer"
)

func TestParseTimeRange(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    indexer.TimeRange
		wantErr bool
	}{
		{
			"Year",
			args{"2015"},
			indexer.TimeRange{time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), indexer.Year},
			false,
		},
		{
			"Month",
			args{"2015-08"},
			indexer.TimeRange{time.Date(2015, 8, 1, 0, 0, 0, 0, time.UTC), indexer.Month},
			false,
		},
		{
			"Day",
			args{"2015-08-03"},
			indexer.TimeRange{time.Date(2015, 8, 3, 0, 0, 0, 0, time.UTC), indexer.Day},
			false,
		},
		{
			"Hour",
			args{"2015-08-01 15"},
			indexer.TimeRange{time.Date(2015, 8, 1, 15, 0, 0, 0, time.UTC), indexer.Hour},
			false,
		},
		{
			"Minute",
			args{"2015-08-01 00:04"},
			indexer.TimeRange{time.Date(2015, 8, 1, 0, 4, 0, 0, time.UTC), indexer.Minute},
			false,
		},
		{
			"ErrorFormat",
			args{"2015-08-01T00:04"},
			indexer.TimeRange{},
			true,
		},
		{
			"NotSupportedFormat",
			args{"2015-08-01 00:04:30"},
			indexer.TimeRange{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := indexer.ParseTimeRange(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeRange_String(t *testing.T) {
	tests := []struct {
		name      string
		timeRange indexer.TimeRange
		want      string
	}{
		{
			"Year",
			indexer.TimeRange{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), indexer.Year},
			"2006",
		},
		{
			"Month",
			indexer.TimeRange{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), indexer.Month},
			"2006-01",
		},
		{
			"Day",
			indexer.TimeRange{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), indexer.Day},
			"2006-01-02",
		},
		{
			"Hour",
			indexer.TimeRange{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), indexer.Hour},
			"2006-01-02 15",
		},
		{
			"Minute",
			indexer.TimeRange{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), indexer.Minute},
			"2006-01-02 15:04",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.timeRange.String(); got != tt.want {
				t.Errorf("TimeRange.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
