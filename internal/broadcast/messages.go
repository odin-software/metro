package broadcast

type ADMessage[T any] struct {
	StationID int64
	Train     T
}
