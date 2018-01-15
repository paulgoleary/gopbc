package field

import "math/big"

type QuadraticField struct {
	BaseField
	targetField *ZField
}

// D2ExtensionQuadElement

var _ PointElement = (*D2ExtensionQuadElement)(nil)
var _ PowElement = (*D2ExtensionQuadElement)(nil)

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

// e0 = (x + y) * (x - y)
// e1 = (xy) * 2
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

func (elem *D2ExtensionQuadElement) Pow(in *ModInt) PointElement {
	result := powWindow(elem, &in.v).(*D2ExtensionQuadElement)
	result.freeze()
	return result
}

func (elem *D2ExtensionQuadElement) CopyPow() PowElement {
	theCopy := elem.dup()
	theCopy.freeze()
	return theCopy
}

func (elem *D2ExtensionQuadElement) MakeOnePow() PowElement {
	return elem.ElemField.makeOneInternal()
}

func (elem *D2ExtensionQuadElement) MulPow(elemIn PowElement) PowElement {
	res := elem.MulPoint(elemIn.(*D2ExtensionQuadElement))
	return res.(*D2ExtensionQuadElement)
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

func (qfield *D2ExtensionQuadField) MakeElementFromBytes(elemBytes []byte) PointElement {
	pnt := MakePointFromBytes(elemBytes, &qfield.targetField.BaseField)
	return qfield.MakeElement(pnt.dataX, pnt.dataY)
}

func MakeD2ExtensionQuadField(Fq *ZField) *D2ExtensionQuadField {

	qfield := new(D2ExtensionQuadField)
	qfield.targetField = Fq
	qfield.FieldOrder = new(big.Int)
	qfield.FieldOrder.Mul(qfield.targetField.FieldOrder, qfield.targetField.FieldOrder)
	qfield.LengthInBytes = qfield.targetField.LengthInBytes * 2

	return qfield
}

func (qfield *D2ExtensionQuadField) makeOneInternal() *D2ExtensionQuadElement {
	return &D2ExtensionQuadElement{qfield,
		PointLike{MakeModInt(1, true, qfield.targetField.FieldOrder), MakeModInt(0, true, qfield.targetField.FieldOrder)}}
}

func (qfield *D2ExtensionQuadField) MakeOne() PointElement {
	return qfield.makeOneInternal()
}
