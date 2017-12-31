package field

import (
	"math/big"
)

type ZrField struct {
	BaseField
	TwoInverse *ModInt
}

type ZrElement struct {
	ElemField *ZrField
	Data      *ModInt
}

// ZrElement

// validate that ZrElement satisfies Element
var _ Element = (*ZrElement)(nil)

func (elem ZrElement) GetInt() *big.Int {
	return &elem.Data.v
}

func (elem ZrElement) String() string {
	return elem.Data.String()
}

// TODO
func (elem ZrElement) Copy() Element {
	return elem
}

// TODO
func (elem ZrElement) Mul(Element) Element {
	return elem
}

// ZrField

func MakeZrField(fieldOrder *big.Int) *ZrField {
	zrField := new(ZrField)
	zrField.FieldOrder = fieldOrder
	zrField.LengthInBytes = fieldOrder.BitLen() / 8 // TODO: generalize ???
	zrField.TwoInverse = zrField.NewElement(TWO).Data.Invert()
	zrField.TwoInverse.Freeze()
	return zrField
}

func (field *ZrField) NewOneElement() *ZrElement {
	return field.NewElement(ONE)
}

func (field *ZrField) NewZeroElement() *ZrElement {
	return field.NewElement(ZERO)
}

func (field *ZrField) NewElement(elemValue *big.Int) *ZrElement {
	elem := new(ZrElement)
	elem.ElemField = field
	elem.Data = CopyFrom(elemValue, true, field.FieldOrder)
	return elem
}
