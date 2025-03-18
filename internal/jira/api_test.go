package jira

import (
	"reflect"
	"testing"
	"time"
)

func TestSprints_GetLast(t *testing.T) {
	tests := []struct {
		name string
		s    Sprints
		want *Sprint
	}{
		{
			name: "Return the last sprint sorted by 'EndDate' field",
			s: []Sprint{
				{
					ID:      1,
					EndDate: time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:      2,
					EndDate: time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &Sprint{
				ID:      2,
				EndDate: time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Return the only sprint p sprint sorted by 'EndDate' field",
			s: []Sprint{
				{
					ID:      1,
					EndDate: time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &Sprint{
				ID:      1,
				EndDate: time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Return nil if the sprint list is empty",
			s:    []Sprint{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetLast(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLast() = %v, want %v", got, tt.want)
			}
		})
	}
}
