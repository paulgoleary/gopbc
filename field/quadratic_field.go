package field

import "math/big"

var _ PointField = (*QuadraticField)(nil)

type QuadraticField struct {
	BaseField
	targetField *ZrField
}

func (field *QuadraticField) MakeElement() PointElement {
	return &QuadraticElement{field, PointLike {nil, nil}}
}

// QuadraticElement

var _ Element = (*QuadraticElement)(nil)
var _ PointElement = (*QuadraticElement)(nil)

type QuadraticElement struct {
	ElemField *QuadraticField
	PointLike
}

func (elem *QuadraticElement) X() *BigInt {
	return elem.DataX
}

func (elem *QuadraticElement) Y() *BigInt {
	return elem.DataY
}

// TODO !!!
func (elem *QuadraticElement) Negate() PointElement {
	return elem
}

func (elem *QuadraticElement) Copy() Element {
	theCopy := elem.dup()
	theCopy.freeze()
	return theCopy
}

func (elem *QuadraticElement) dup() *QuadraticElement {
	newElem := new(QuadraticElement)
	newElem.ElemField = elem.ElemField
	newElem.DataX = elem.DataX.copy()
	newElem.DataY = elem.DataY.copy()
	return newElem
}

func (elem *QuadraticElement) SetToOne() Element {
	return &QuadraticElement{elem.ElemField, PointLike{BI_ONE, BI_ZERO}}
}

func (elem *QuadraticElement) Mul(Element) Element {
	return nil // TODO!!!
}

// D2ExtensionQuadField

type D2ExtensionQuadField struct {
	QuadraticField
}

func MakeD2ExtensionQuadField(Fq *ZrField) *D2ExtensionQuadField {

	field := new(D2ExtensionQuadField)
	field.targetField = Fq
	field.FieldOrder = new(big.Int)
	field.FieldOrder.Mul(field.targetField.FieldOrder, field.targetField.FieldOrder)
	field.LengthInBytes = field.targetField.LengthInBytes * 2

	return field
}