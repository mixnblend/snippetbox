package main

import (
	"testing"
	"time"

	"github.com/mixnblend/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name     string
		datetime time.Time
		want     string
	}{
		{
			name:     "UTC",
			datetime: time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want:     "17 Mar 2024 at 10:15",
		},
		{
			name:     "Empty",
			datetime: time.Time{},
			want:     "",
		},
		{
			name:     "CET",
			datetime: time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want:     "17 Mar 2024 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given ... we have a date
			datetime := tt.datetime

			// When ... we call our function
			result := humanDate(datetime)

			// Then ... the date should be formatted in a human readable form as expected
			assert.Equal(t, result, tt.want)
		})
	}
}
