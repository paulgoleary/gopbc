package field

import (
	"math/big"
	"log"
	"fmt"
)

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

/*
func (elem *D2ExtensionQuadElement) MulPoint(elemIn PointElement) PointElement {
	e2 := elem.dataX.Add(elem.dataY).Mul(elemIn.X().Add(elemIn.Y()))
	e0 := elem.dataX.Mul(elemIn.X())
	e1 := elem.dataY.Mul(elemIn.Y())
	e2 = e2.Sub(e0)
	return elem.ElemField.MakeElement(e0.Sub(e1), e2.Sub(e1))
}
*/

func (elem *D2ExtensionQuadElement) MulPoint(elemIn PointElement) PointElement {
	return elem.ElemField.ElemMul(elem, elemIn)
}

// TODO: move elsewhere ...?
func (qfield *D2ExtensionQuadField) ElemMul(elemA PointElement, elemB PointElement) PointElement {
	e2 := elemA.X().Add(elemA.Y()).Mul(elemB.X().Add(elemB.Y()))
	e0 := elemA.X().Mul(elemB.X())
	e1 := elemA.Y().Mul(elemB.Y())
	e2 = e2.Sub(e0)
	return qfield.MakeElement(e0.Sub(e1), e2.Sub(e1))
}

func (qfield *D2ExtensionQuadField) ElemMulNor(elem PointElement) PointElement {
	// TODO: don't necessarily need this check. but for now...
	testMod8 := new(big.Int).Mod(qfield.FieldOrder, big.NewInt(8))
	if testMod8.Cmp(THREE) != 0 {
		log.Panicf("Currently only implemented for field order 3 % 8")
	}
	t0 := elem.Y().Negate()
	c1 := elem.X().Add(elem.Y())
	c0 := t0.Add(elem.X())
	return qfield.MakeElement(c0, c1)
}

func (qfield *D2ExtensionQuadField) ElemAdd(elemA PointElement, elemB PointElement) PointElement {
	return qfield.MakeElement(elemA.X().Add(elemB.X()), elemB.X().Add(elemB.Y()))
}

func (qfield *D2ExtensionQuadField) ElemSub(elemA PointElement, elemB PointElement) PointElement {
	return qfield.MakeElement(elemA.X().Sub(elemB.X()), elemB.X().Sub(elemB.Y()))
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

// D6ExtensionQuadElement

type D6ExtensionQuadField struct {
	BaseField
	targetField *D2ExtensionQuadField
}

func MakeD6ExtensionQuadField(d2Ext *D2ExtensionQuadField) *D6ExtensionQuadField {

	qfield := new(D6ExtensionQuadField)
	qfield.targetField = d2Ext
	qfield.FieldOrder = new(big.Int)
	qfield.FieldOrder.Exp(qfield.targetField.FieldOrder, THREE, nil) // not so sure about this ... :/
	qfield.LengthInBytes = qfield.targetField.LengthInBytes * 3

	return qfield
}

// An element of Fq6, represented by c0 + c1 * v + c2 * v^(2).
type D6ExtensionQuadElement struct {
	ElemField *D6ExtensionQuadField
	c0 PointElement
	c1 PointElement
	c2 PointElement
}

func (field *D6ExtensionQuadField) MakeElement(c0, c1, c2 PointElement) *D6ExtensionQuadElement {
	elem := new(D6ExtensionQuadElement)
	elem.ElemField = field
	elem.c0 = c0
	elem.c1 = c1
	elem.c2 = c2
	return elem
}

func (elem *D6ExtensionQuadElement) String() string {
	return fmt.Sprintf("D6 elem: [%s,\n%s],\n[%s,\n%s],\n[%s,\n%s]",
		elem.c0.X().String(), elem.c0.Y().String(),
		elem.c1.X().String(), elem.c1.Y().String(),
		elem.c2.X().String(), elem.c2.Y().String())
}

func (elem *D6ExtensionQuadElement) MulPoint(elemIn *D6ExtensionQuadElement) *D6ExtensionQuadElement {

	targetField := elem.ElemField.targetField

	// v0 = a_0 * b_0
	v0 := targetField.ElemMul(elem.c0, elemIn.c0)
	Trace(v0)

	// v1 = a_1 * b_1
	v1 := targetField.ElemMul(elem.c1, elemIn.c1)
	Trace(v1)

	// v2 = a_2 * b_2
	v2 := targetField.ElemMul(elem.c2, elemIn.c2)
	Trace(v2)

	// t2 (c_0) = v0 + E((a_1 + a_2)(b_1 + b_2) - v1 - v2)
	// (a_1 + a_2)
	t0 := targetField.ElemAdd(elem.c1, elem.c2)
	// (b_1 + b_2)
	t1 := targetField.ElemAdd(elemIn.c1, elemIn.c2)
	// (a_1 + a_2)(b_1 + b_2)
	t2 := targetField.ElemMul(t0, t1)

	t2 = targetField.ElemSub(t2, v1)
	t2 = targetField.ElemSub(t2, v2)

	// fp2_mul_nor(t0, t2);
	t0 = targetField.ElemMulNor(t2)
	// fp2_add(t2, t0, v0);
	t2 = targetField.ElemAdd(t0, v0)

	return nil
}
