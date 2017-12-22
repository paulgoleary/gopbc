package field

import (
	"math/big"
)

type ZrField struct {
	BaseField
}

type ZrElement struct {
	ElemField *ZrField
	Data      *BigInt
}

// ZrElement

// validate that ZrElement satisfies Element
var _ Element = (*ZrElement)(nil)

// TODO
func (elem ZrElement) PowZn(eZn ZrElement) ZrElement {
	return elem
}

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

// TODO
func (elem ZrElement) SetToOne() Element {
	return elem
}

// TODO
func (elem ZrElement) Square() Element {
	return elem
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

func (field *ZrField) NewZero() *BigInt {
	return BI_ZERO // ok to use because it's frozen
}

func (field *ZrField) NewElement(elemValue *big.Int) *ZrElement {
	elem := new(ZrElement)
	elem.ElemField = field
	elem.Data = CopyFrom(elemValue, true)
	return elem
}
