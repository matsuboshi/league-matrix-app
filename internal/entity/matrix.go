package entity

// Matrix represents a two-dimensional matrix of integer values.
// The Data field contains rows of integer columns, where each row must have the same length.
type Matrix struct {
	Data [][]int64
}
