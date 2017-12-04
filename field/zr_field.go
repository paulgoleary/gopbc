package field

import (
	"math/big"
)

type ZrField struct {
	FieldOrder	*big.Int
}

type ZrElement struct {
	ElemField *ZrField
	Data      *big.Int
}

// ZrElement

func (e ZrElement) PowZn(eZn ZrElement) ZrElement {
	return e
}

// ZrField

func MakeZrField( fieldOrder *big.Int ) *ZrField {
	zrField := new(ZrField)
	zrField.FieldOrder = fieldOrder
	return zrField
}

func (field *ZrField) NewOneElement() *ZrElement {
	elem := new(ZrElement)
	elem.ElemField = field
	elem.Data = ONE
	return elem
}

func (field *ZrField) NewZeroElement() *ZrElement {
	elem := new(ZrElement)
	elem.ElemField = field
	elem.Data = ZERO
	return elem
}