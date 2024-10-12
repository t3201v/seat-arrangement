package model

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestCinema_IsValidGroup(t *testing.T) {
	type fields struct {
		logger      *log.Logger
		rows        int
		columns     int
		minDistance int
		seats       [][]Seat
	}
	type args struct {
		seatCoords [][]int
		groupName  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "pass /w same group",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]Seat{
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
				},
			},
			args: args{
				seatCoords: [][]int{{0, 0}, {3, 4}},
				groupName:  "pass",
			},
			want: true,
		},
		{
			name: "fail bc different groups and not met min distance",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]Seat{
					{{1, "a"}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
				},
			},
			args: args{
				seatCoords: [][]int{{3, 4}},
				groupName:  "b",
			},
			want: false,
		},
		{
			name: "fail bc buying reserved seat",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]Seat{
					{{1, "a"}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {0, ""}},
					{{0, ""}, {0, ""}, {0, ""}, {0, ""}, {1, ""}},
				},
			},
			args: args{
				seatCoords: [][]int{{3, 4}},
				groupName:  "a",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cinema{
				logger:      tt.fields.logger,
				rows:        tt.fields.rows,
				columns:     tt.fields.columns,
				minDistance: tt.fields.minDistance,
				seats:       tt.fields.seats,
			}
			if got := c.IsValidGroup(tt.args.seatCoords, tt.args.groupName); got != tt.want {
				t.Errorf("IsValidGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
