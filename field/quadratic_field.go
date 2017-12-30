package field

import "math/big"

type QuadraticField struct {
	BaseField
	targetField *ZrField
}

// D2ExtensionQuadElement

var _ PointElement = (*D2ExtensionQuadElement)(nil)

type D2ExtensionQuadElement struct {
	ElemField *D2ExtensionQuadField
	PointLike
}

func (elem *D2ExtensionQuadElement) X() *BigInt {
	return elem.dataX
}

func (elem *D2ExtensionQuadElement) Y() *BigInt {
	return elem.dataY
}

// TODO: logic is very similar to curve field ...?
func (elem *D2ExtensionQuadElement) NegateY() PointElement {
	if elem.dataY.IsEqual(BI_ZERO) {
		return elem
	}
	elem.PointLike.freeze() // make sure we're frozen
	yNeg := elem.dataY.Negate(elem.ElemField.targetField.FieldOrder)
	return &D2ExtensionQuadElement{elem.ElemField, PointLike{elem.dataX, yNeg}}
}

func (elem *D2ExtensionQuadElement) Invert() PointElement {

	targetOrder := elem.ElemField.targetField.FieldOrder
	elem.freeze()

	e0:= elem.dataX.Square(targetOrder).Add(elem.dataY.Square(targetOrder), targetOrder).Invert(targetOrder)

	x := elem.dataX.Mul(e0, targetOrder)
	// no need to freeze (and therefore Copy) e0 because it's not used again
	y := elem.dataY.Mul(e0.Negate(targetOrder), targetOrder)

	return elem.ElemField.MakeElement(x, y)
}

func (elem *D2ExtensionQuadElement) dup() *D2ExtensionQuadElement {
	newElem := new(D2ExtensionQuadElement)
	newElem.ElemField = elem.ElemField
	newElem.dataX = elem.dataX.Copy()
	newElem.dataY = elem.dataY.Copy()
	return newElem
}

func (elem *D2ExtensionQuadElement) Square() PointElement {
	targetOrder := elem.ElemField.targetField.FieldOrder // TODO: verify
	e0 := elem.dataX.Add(elem.dataY, targetOrder).Mul(elem.dataX.Sub(elem.dataY, targetOrder), targetOrder)
	e1 := elem.dataX.Mul(elem.dataY, targetOrder).Mul(BI_TWO, targetOrder)
	return elem.ElemField.MakeElement(e0, e1)
}

func (elem *D2ExtensionQuadElement) MulPoint(elemIn PointElement) PointElement {
	targetOrder := elem.ElemField.targetField.FieldOrder // TODO - verify !
	e2 := elem.dataX.Add(elem.dataY, targetOrder).Mul(elemIn.X().Add(elemIn.Y(), targetOrder), targetOrder)
	e0 := elem.dataX.Mul(elemIn.X(), targetOrder)
	e1 := elem.dataY.Mul(elemIn.Y(), targetOrder)
	e2 = e2.Sub(e0, targetOrder)
	return elem.ElemField.MakeElement(e0.Sub(e1, targetOrder), e2.Sub(e1, targetOrder))
}

// D2ExtensionQuadField

var _ PointField = (*D2ExtensionQuadField)(nil)

type D2ExtensionQuadField struct {
	QuadraticField
}

func (field *D2ExtensionQuadField) MakeElement(x *BigInt, y *BigInt) PointElement {
	if x != nil {
		x.Freeze()
	}
	if y != nil {
		y.Freeze()
	}
	return &D2ExtensionQuadElement{field, PointLike {x, y}}
}

func MakeD2ExtensionQuadField(Fq *ZrField) *D2ExtensionQuadField {

	field := new(D2ExtensionQuadField)
	field.targetField = Fq
	field.FieldOrder = new(big.Int)
	field.FieldOrder.Mul(field.targetField.FieldOrder, field.targetField.FieldOrder)
	field.LengthInBytes = field.targetField.LengthInBytes * 2

	return field
}

func (field *D2ExtensionQuadField) MakeOne() PointElement {
	return &D2ExtensionQuadElement{field, PointLike{BI_ONE, BI_ZERO}}
}
