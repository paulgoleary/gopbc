package field

import (
	"math/big"
)

type ZrField struct {
	BaseField
	FieldOrder *big.Int
}

type ZrElement struct {
	ElemField *ZrField
	Data      *BigInt
}

// ZrElement

func (e ZrElement) PowZn(eZn ZrElement) ZrElement {
	return e
}

// ZrField

func MakeZrField(fieldOrder *big.Int) *ZrField {
	zrField := new(ZrField)
	zrField.FieldOrder = fieldOrder
	zrField.LengthInBytes = fieldOrder.BitLen() / 8 // TODO: generalize ???
	return zrField
}

func (field *ZrField) NewOneElement() *ZrElement {
	elem := new(ZrElement)
	elem.ElemField = field
	elem.Data = BI_ONE
	return elem
}

func (field *ZrField) NewZeroElement() *ZrElement {
	elem := new(ZrElement)
	elem.ElemField = field
	elem.Data = BI_ZERO
	return elem
}
