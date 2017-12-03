package field

import (
	"math/big"
)

type ZrField struct {
	FieldOrder	big.Int
}

type ZrElement struct {
	ElemField ZrField
	Data      big.Int
}

// ZrElement

func (e ZrElement) PowZn(eZn ZrElement) ZrElement {
	return e
}

// ZrField

func MakeZrField( fieldOrder *big.Int ) *ZrField {
	zrField := new(ZrField)
	zrField.FieldOrder = *fieldOrder
	return zrField
}