package kernel

import (
	"fmt"
	"math/rand"
)
const COORDINATE_MIN_VALUE = 1
const COORDINATE_MAX_VALUE = 10

type Location struct {
	x        uint8
	y        uint8
	isSet bool
}

func NewLocation(x uint8, y uint8) (Location, error) {
	if x > COORDINATE_MAX_VALUE || x < COORDINATE_MIN_VALUE {
		return Location{ isSet: false }, fmt.Errorf("переданы невалидные значения координаты x: %d", x)
	}

	if y > COORDINATE_MAX_VALUE || y < COORDINATE_MIN_VALUE {
		return Location{ isSet: false }, fmt.Errorf("переданы невалидные значения координаты y: %d", y)
	}

	location := Location{
		x:     x,
		y:     y,
		isSet: true,
	}

	return location, nil
}

func NewRandomLocation() (*Location, error) {
	location := Location{
		x: uint8(rand.Intn(COORDINATE_MAX_VALUE) + 1),
		y: uint8(rand.Intn(COORDINATE_MAX_VALUE) + 1),
		isSet: true,
	}

	return &location, nil
}

func (l Location) X() uint8 {
	return l.x;
}

func (l Location) Y() uint8 {
	return l.y;
}

func (l Location) Equals(other Location) bool {
	return l == other;
}

func (l Location) IsEmpty() bool {
	return l.isSet == false
}

func (l Location) DistanceTo(target Location) (Distance, error) {
	return NewDistance(l, target)
}
