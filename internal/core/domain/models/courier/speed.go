package courier

import "errors"

type Speed struct {
	value int
}

func NewSpeed(value int) (Speed, error) {
	if value <= 0 {
		return Speed{}, errors.New("скорость курьера должна быть больше нуля")
	}
	return Speed{value: value}, nil
}

func (s Speed) Value() int {
	return s.value
}

func (s Speed) Equals(other Speed) bool {
	return s.value == other.value
}
