package broadcast

type ADMessage[T any] struct {
	StationID string
	Train     T
}
