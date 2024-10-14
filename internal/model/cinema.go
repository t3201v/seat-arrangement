package model

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/t3201v/seat-arrangement/gen/cinema"
	"github.com/t3201v/seat-arrangement/internal/helper"
)

type SeatStatus int

const (
	Available SeatStatus = iota
	Reserved
	SeatStatusEnd
)

const (
	Row int = iota
	Col
)

var directions = [][]int{
	{-1, 0}, // up
	{1, 0},  // down
	{0, -1}, // left
	{0, 1},  // right
}

type Seat struct {
	status    SeatStatus
	groupName string
}

// Cinema structure to hold seat information and minimum distance rule
type Cinema struct {
	logger      *log.Logger
	rows        int
	columns     int
	minDistance int
	seats       [][]Seat // 0: available, 1: reserved
}

// NewCinema initializes the cinema layout with the given rows, columns, and min_distance
func NewCinema(l *log.Logger, rows, columns, minDistance int) *Cinema {
	seats := make([][]Seat, rows)
	for i := range seats {
		seats[i] = make([]Seat, columns)
	}
	return &Cinema{
		logger:      l,
		rows:        rows,
		columns:     columns,
		minDistance: minDistance,
		seats:       seats,
	}
}

func (c *Cinema) validate(seatCoords [][]int) error {
	if len(c.seats) == 0 || len(c.seats[0]) == 0 {
		return errors.New("malformed seats data")
	}
	for _, seat := range seatCoords {
		if len(seat) != 2 {
			return errors.New("seat coordinates must be of length 2")
		}
		if seat[Row] < 0 || seat[Row] >= len(c.seats) {
			return fmt.Errorf("seat coordinates must be in range [%d, %d)", Row, c.rows)
		}
		if seat[Col] < 0 || seat[Col] >= len(c.seats[Row]) {
			return fmt.Errorf("seat coordinates must be in range [%d, %d)", Col, c.columns)
		}
	}
	return nil
}

// it will reset seats if rows or columns number's changed
func (c *Cinema) UpdateConfig(rows, columns, minDistance int) {
	if c.rows != rows || c.columns != columns {
		seats := make([][]SeatStatus, rows)
		for i := range seats {
			seats[i] = make([]SeatStatus, columns)
		}
	}
	c.rows = rows
	c.columns = columns
	c.minDistance = minDistance
}

// IsValidGroup checks if a group of seats can be reserved together
func (c *Cinema) IsValidGroup(seatCoords [][]int, groupName string) bool {
	if err := c.validate(seatCoords); err != nil {
		return false
	}
	if len(seatCoords) == 0 {
		return false
	}

	// Check if the group contains any reserved seats
	for _, seat := range seatCoords {
		if c.seats[seat[Row]][seat[Col]].status == Reserved {
			return false // Seat is already reserved
		}
	}

	// Check if the group satisfies the minimum distance rule
	for i := 0; i < c.rows; i++ {
		for j := 0; j < c.columns; j++ {
			if c.seats[i][j].status != Reserved || c.seats[i][j].groupName == groupName {
				continue
			}
			// check minimum distance for different groups of seats
			for _, seat := range seatCoords {
				if helper.ManhattanDistance(i, j, seat[Row], seat[Col]) <= c.minDistance {
					return false
				}
			}
		}
	}
	return true
}

// ReserveSeats attempts to reserve seats if they are valid according to the distance rule
func (c *Cinema) ReserveSeats(seatCoords [][]int, groupName string) error {
	if err := c.validate(seatCoords); err != nil {
		return err
	}

	if !c.IsValidGroup(seatCoords, groupName) {
		return errors.New("seats are not available right now")
	}
	for _, seat := range seatCoords {
		c.seats[seat[Row]][seat[Col]].status = Reserved
		c.seats[seat[Row]][seat[Col]].groupName = groupName
	}
	return nil
}

// CancelSeats cancels the reservation of specific seats
func (c *Cinema) CancelSeats(seatCoords [][]int) error {
	if err := c.validate(seatCoords); err != nil {
		return err
	}

	for _, seat := range seatCoords {
		if c.seats[seat[Row]][seat[Col]].status == Available {
			return fmt.Errorf("seat (%d, %d) is not reserved", seat[Row], seat[Col])
		}
		c.seats[seat[Row]][seat[Col]].status = Available
		c.seats[seat[Row]][seat[Col]].groupName = ""
	}
	return nil
}

// ListAvailableSeatsGrouped returns groups of available seats that can be purchased together
func (c *Cinema) ListAvailableSeatsGrouped() [][]int {
	availableGroups := make([][]int, 0)
	for i := 0; i < c.rows; i++ {
		for j := 0; j < c.columns; j++ {
			if c.seats[i][j].status == Available {
				availableGroups = append(availableGroups, []int{i, j})
			}
		}
	}

	return availableGroups
}

func (c *Cinema) ToPbSeats(seats [][]int) ([]*cinema.Seat, error) {
	result := make([]*cinema.Seat, 0)
	for _, seat := range seats {
		if len(seat) < 2 {
			return nil, errors.New("malformed seat data")
		}
		result = append(result, &cinema.Seat{
			Row:    int32(seat[Row]),
			Column: int32(seat[Col]),
		})
	}
	return result, nil
}

// PrintLayout prints the current layout of the cinema (for testing purposes)
func (c *Cinema) String() string {
	var sb strings.Builder

	for i, row := range c.seats {
		for j, val := range row {
			sb.WriteString(fmt.Sprintf("%d", val.status))
			if j < len(row)-1 {
				sb.WriteString(" ") // Add a space between numbers in the same row
			}
		}
		if i < len(c.seats)-1 {
			sb.WriteString("\n") // Add a newline between rows
		}
	}

	return sb.String()
}

func (c *Cinema) Clone() *Cinema {
	// Create a new Cinema struct
	newCinema := &Cinema{
		logger:      c.logger,
		rows:        c.rows,
		columns:     c.columns,
		minDistance: c.minDistance,
		seats:       make([][]Seat, len(c.seats)),
	}

	// Deep copy of the seats slice
	for i := range c.seats {
		newCinema.seats[i] = make([]Seat, len(c.seats[i]))
		copy(newCinema.seats[i], c.seats[i]) // Copy the seat statuses row by row
	}

	return newCinema
}
