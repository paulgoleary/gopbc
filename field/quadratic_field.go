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

func (elem *D2ExtensionQuadElement) X() *ModInt {
	return elem.dataX
}

func (elem *D2ExtensionQuadElement) Y() *ModInt {
	return elem.dataY
}

// TODO: logic is very similar to curve field ...?
func (elem *D2ExtensionQuadElement) NegateY() PointElement {
	if elem.dataY.IsValEqual(MI_ZERO) {
		return elem
	}
	elem.PointLike.freeze() // make sure we're frozen
	yNeg := elem.dataY.Negate()
	return &D2ExtensionQuadElement{elem.ElemField, PointLike{elem.dataX, yNeg}}
}

func (elem *D2ExtensionQuadElement) Invert() PointElement {

	elem.freeze()

	e0:= elem.dataX.Square().Add(elem.dataY.Square()).Invert()

	x := elem.dataX.Mul(e0)
	// no need to freeze (and therefore Copy) e0 because it's not used again
	y := elem.dataY.Mul(e0.Negate())

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
	e0 := elem.dataX.Add(elem.dataY).Mul(elem.dataX.Sub(elem.dataY))
	e1 := elem.dataX.Mul(elem.dataY).Mul(MI_TWO)
	return elem.ElemField.MakeElement(e0, e1)
}

func (elem *D2ExtensionQuadElement) MulPoint(elemIn PointElement) PointElement {
	e2 := elem.dataX.Add(elem.dataY).Mul(elemIn.X().Add(elemIn.Y()))
	e0 := elem.dataX.Mul(elemIn.X())
	e1 := elem.dataY.Mul(elemIn.Y())
	e2 = e2.Sub(e0)
	return elem.ElemField.MakeElement(e0.Sub(e1), e2.Sub(e1))
}

// D2ExtensionQuadField

var _ PointField = (*D2ExtensionQuadField)(nil)

type D2ExtensionQuadField struct {
	QuadraticField
}

func (qfield *D2ExtensionQuadField) MakeElement(x *ModInt, y *ModInt) PointElement {
	if x != nil {
		x.Freeze()
	}
	if y != nil {
		y.Freeze()
	}
	return &D2ExtensionQuadElement{qfield, PointLike {x, y}}
}

func MakeD2ExtensionQuadField(Fq *ZrField) *D2ExtensionQuadField {

	qfield := new(D2ExtensionQuadField)
	qfield.targetField = Fq
	qfield.FieldOrder = new(big.Int)
	qfield.FieldOrder.Mul(qfield.targetField.FieldOrder, qfield.targetField.FieldOrder)
	qfield.LengthInBytes = qfield.targetField.LengthInBytes * 2

	return qfield
}

func (qfield *D2ExtensionQuadField) MakeOne() PointElement {
	return &D2ExtensionQuadElement{qfield,
	PointLike{MakeModInt(1, true, qfield.targetField.FieldOrder), MakeModInt(0, true, qfield.targetField.FieldOrder)}}
}
