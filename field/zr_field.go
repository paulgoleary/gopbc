package field

import (
	"math/big"
	"fmt"
)

type ZrField struct {
	BaseField
	TwoInverse *ModInt
}

type ZrElement struct {
	ElemField *ZrField
	*ModInt
}

func (elem *ZrElement) String() string {
	return fmt.Sprintf("ZrElement: %s | %s", elem.v.String(), elem.ElemField.FieldOrder.String())
}

// ZrField

func MakeZrField(fieldOrder *big.Int) *ZrField {
	zrField := new(ZrField)
	zrField.FieldOrder = fieldOrder
	zrField.LengthInBytes = fieldOrder.BitLen() / 8 // TODO: generalize ???
	zrField.TwoInverse = zrField.NewElement(TWO).Invert()
	zrField.TwoInverse.Freeze()
	return zrField
}

func (zrfield *ZrField) NewOneElement() *ZrElement {
	return zrfield.NewElement(ONE)
}

func (zrfield *ZrField) NewZeroElement() *ZrElement {
	return zrfield.NewElement(ZERO)
}

func (zrfield *ZrField) NewElement(elemValue *big.Int) *ZrElement {
	return &ZrElement{zrfield,CopyFrom(elemValue, true, zrfield.FieldOrder)}
}