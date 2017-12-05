package field

import "math/big"

var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var TWO = big.NewInt(2)

type Field interface {
}

// TODO: not sure if I want/like/need this ...?
type Element interface {
	Copy() Element
}

type BaseField struct {
	LengthInBytes int
}

type Point interface {
	X() *big.Int
	Y() *big.Int
}

