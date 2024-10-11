package helper

import (
	"github.com/t3201v/seat-arrangement/internal/libs/util"
)

// ManhattanDistance calculates the Manhattan distance between two points
func ManhattanDistance(x1, y1, x2, y2 int) int {
	return util.Abs(x1-x2) + util.Abs(y1-y2)
}
