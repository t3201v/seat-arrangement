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

// Cinema structure to hold seat information and minimum distance rule
type Cinema struct {
	logger      *log.Logger
	rows        int
	columns     int
	minDistance int
	seats       [][]SeatStatus // 0: available, 1: reserved
}

// NewCinema initializes the cinema layout with the given rows, columns, and min_distance
func NewCinema(l *log.Logger, rows, columns, minDistance int) *Cinema {
	seats := make([][]SeatStatus, rows)
	for i := range seats {
		seats[i] = make([]SeatStatus, columns)
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

func (c *Cinema) UpdateConfig(rows, columns, minDistance int) {
	c.rows = rows
	c.columns = columns
	c.minDistance = minDistance
}

// IsValidGroup checks if a group of seats can be reserved together
func (c *Cinema) IsValidGroup(seatCoords [][]int) bool {
	if err := c.validate(seatCoords); err != nil {
		return false
	}
	if len(seatCoords) == 0 {
		return false
	}

	// Check all pairs of seats in the group
	for i := 0; i < len(seatCoords); i++ {
		for j := i + 1; j < len(seatCoords); j++ {
			// Calculate Manhattan distance
			distance := helper.ManhattanDistance(seatCoords[i][Row], seatCoords[i][Col], seatCoords[j][Row], seatCoords[j][Col])
			if distance <= c.minDistance {
				return false // Not satisfying the minimum distance rule
			}
		}
	}

	// Check if the group contains any reserved seats
	for _, seat := range seatCoords {
		if c.seats[seat[Row]][seat[Col]] == Reserved {
			return false // Seat is already reserved
		}
	}

	// Check if the group satisfies the minimum distance rule
	for i := 0; i < c.rows; i++ {
		for j := 0; j < c.columns; j++ {
			if c.seats[i][j] == Reserved { // Already reserved seat or to be reserved seat
				for _, seat := range seatCoords {
					if helper.ManhattanDistance(i, j, seat[Row], seat[Col]) <= c.minDistance {
						return false
					}
				}
			}
		}
	}
	return true
}

// ReserveSeats attempts to reserve seats if they are valid according to the distance rule
func (c *Cinema) ReserveSeats(seatCoords [][]int) error {
	if err := c.validate(seatCoords); err != nil {
		return err
	}

	if !c.IsValidGroup(seatCoords) {
		return errors.New("seats cannot be reserved due to minimum distance rule")
	}
	for _, seat := range seatCoords {
		c.seats[seat[Row]][seat[Col]] = Reserved
	}
	return nil
}

// CancelSeats cancels the reservation of specific seats
func (c *Cinema) CancelSeats(seatCoords [][]int) error {
	if err := c.validate(seatCoords); err != nil {
		return err
	}

	for _, seat := range seatCoords {
		if c.seats[seat[Row]][seat[Col]] == Available {
			return fmt.Errorf("seat (%d, %d) is not reserved", seat[Row], seat[Col])
		}
		c.seats[seat[Row]][seat[Col]] = Available
	}
	return nil
}

// ListAvailableSeats returns a list of available seats in the cinema
func (c *Cinema) ListAvailableSeats() [][]int {
	availableSeats := make([][]int, 0)
	for i := 0; i < c.rows; i++ {
		for j := 0; j < c.columns; j++ {
			if c.seats[i][j] == Available {
				availableSeats = append(availableSeats, []int{i, j})
			}
		}
	}
	return availableSeats
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
			sb.WriteString(fmt.Sprintf("%d", val))
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
		seats:       make([][]SeatStatus, len(c.seats)),
	}

	// Deep copy of the seats slice
	for i := range c.seats {
		newCinema.seats[i] = make([]SeatStatus, len(c.seats[i]))
		copy(newCinema.seats[i], c.seats[i]) // Copy the seat statuses row by row
	}

	return newCinema
}
