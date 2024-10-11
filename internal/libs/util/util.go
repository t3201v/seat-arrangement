package util

func Abs[T int | uint | float32 | float64](a T) T {
	if a > 0 {
		return a
	}
	return -a
}
