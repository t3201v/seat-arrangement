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
		seats       [][]SeatStatus
	}
	type args struct {
		seatCoords [][]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "bad arguments",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]SeatStatus{
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
				},
			},
			args: args{
				seatCoords: [][]int{{0, 0}, {3, 4}},
			},
			want: false,
		},
		{
			name: "existed seat and farthest invalid reservation seat",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]SeatStatus{
					{1, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
				},
			},
			args: args{
				seatCoords: [][]int{{3, 4}},
			},
			want: false,
		},
		{
			name: "already reserved",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]SeatStatus{
					{1, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 1},
				},
			},
			args: args{
				seatCoords: [][]int{{3, 4}},
			},
			want: false,
		},
		{
			name: "existed seat and invalid reservation",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 7,
				seats: [][]SeatStatus{
					{0, 0, 0, 0, 0},
					{0, 1, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
				},
			},
			args: args{
				seatCoords: [][]int{{3, 4}},
			},
			want: false,
		},
		{
			name: "pass",
			fields: fields{
				logger:      log.StandardLogger(),
				rows:        4,
				columns:     5,
				minDistance: 6,
				seats: [][]SeatStatus{
					{1, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
				},
			},
			args: args{
				seatCoords: [][]int{{3, 4}},
			},
			want: true,
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
			if got := c.IsValidGroup(tt.args.seatCoords); got != tt.want {
				t.Errorf("IsValidGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
