package kernel

import "fmt"

type Distance struct {
	value int
}

func NewDistance(from Location, to Location) (Distance, error) {
	if !from.isSet || !to.isSet {
		return Distance{}, fmt.Errorf("одна из координат некорректно установлена")
	}

	distanceValue := absInt(int8(from.x-to.x)) + absInt(int8(from.y-to.y))
	return Distance{value: distanceValue}, nil
}

func (d Distance) Value() int {
	return d.value
}

func (d Distance) Equals(other Distance) bool {
	return d.value == other.value
}

func (d Distance) IsZero() bool {
	return d.value == 0
}

func absInt(x int8) int {
	if x < 0 {
		return int(-x)
	}
	return int(x)
}
