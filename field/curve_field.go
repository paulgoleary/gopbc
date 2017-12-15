package field

import (
	"math/big"
	"fmt"
)

type CurveField struct {
	a          *ZrElement
	b          *ZrElement
	order      *big.Int
	cofactor   *big.Int      // TODO: do we need this ...?
	gen        *CurveElement // TODO: not sure here...
	genNoCofac *CurveElement // TODO: don't need this ...?
}

type CurveElement struct {
	ElemField *CurveField
	DataX     *BigInt // TODO: perhaps X and Y should be elements of the target field, as in PBC/JPBC?
	DataY     *BigInt
}

// CurveField

// TODO: JPBC (PBC?) handles case w/o bytes and cofactor
func (field *CurveField) initGenFromBytes(genNoCofac *[]byte) {
	field.genNoCofac = field.newElementFromBytes(genNoCofac)
	field.gen = field.genNoCofac.MulScalar(field.cofactor)
	if !field.gen.isValid(){
		panic("Curve field generator needs to be valid")
	}
}

func (field *CurveField) GetGen() *CurveElement {
	return field.gen
}

func (field *CurveField) getTargetField() *ZrField {
	return field.a.ElemField
}

func (field *CurveField) newElementFromBytes(elemBytes *[]byte) *CurveElement {

	elem := CurveElement{ElemField: field}

	xBytes := (*elemBytes)[:field.getTargetField().LengthInBytes]
	yBytes := (*elemBytes)[field.getTargetField().LengthInBytes:]

	elem.DataX = new(BigInt)
	elem.DataX.setBytes(xBytes)

	elem.DataY = new(BigInt)
	elem.DataY.setBytes(yBytes)

	// needs to be frozen before validation
	elem.freeze()
	if !elem.isValid() {
		elem.setInf()
	}
	return &elem
}

// general curve is y^2 = x^3 + ax + b
func (field *CurveField) calcYSquared(xIn *BigInt) *BigInt {
	order := field.getTargetField().FieldOrder
	if !xIn.frozen {
		panic("xIn needs to be frozen")
	}
	return xIn.square(order).add(field.a.Data, order).mul(xIn, order).add(field.b.Data, order)
}

// TODO: needs to account for sign
func (field *CurveField) newElementFromX(x *big.Int) *CurveElement {

	elem := CurveElement{ElemField: field}

	copyX := CopyFrom(x, true)
	calcY2 := field.calcYSquared(copyX)
	elem.DataY = calcY2.sqrt(field.getTargetField().FieldOrder)
	elem.DataX = copyX

	elem.freeze()
	return &elem
}

func (field *CurveField) newElementFromStrings(xStr string, yStr string) *CurveElement {
	elem := CurveElement{ElemField: field}
	elem.DataX = MakeBigIntStr(xStr, true)
	elem.DataY = MakeBigIntStr(yStr, true)
	return &elem
}

func MakeCurveField(
	a *ZrElement,
	b *ZrElement,
	order *big.Int,
	cofactor *big.Int,
	genNoCofacBytes *[]byte) *CurveField {

	field := new(CurveField)
	field.a = a
	field.b = b
	field.order = order
	field.cofactor = cofactor
	field.initGenFromBytes(genNoCofacBytes)

	return field
}

// make minimal field for testing purposes - TODO: might need a generator?
func makeTestCurveField(a *big.Int, b *big.Int, r *big.Int, q *big.Int) *CurveField {

	zfield := MakeZrField(q)

	cfield := new(CurveField)
	cfield.a = zfield.NewElement(a)
	cfield.b = zfield.NewElement(b)
	cfield.order = r

	return cfield
}

// CurveElement

// TODO: Make function?

// validate that CurveElement satisfies Element
var _ Element = (*CurveElement)(nil)

func (elem *CurveElement) isInf() bool {
	return elem.DataY == nil && elem.DataY == nil
}

func (elem *CurveElement) setInf() {
	elem.DataX = nil
	elem.DataY = nil
}

// satisfy PointLike interface
func (elem *CurveElement) X() *BigInt {
	return elem.DataX
}

func (elem *CurveElement) Y() *BigInt {
	return elem.DataY
}

func (elem *CurveElement) String() string {
	return fmt.Sprintf("[%s],[%s]", elem.DataX.String(), elem.DataY.String())
}

// don't return elem to emmphasize that call mutates elem
func (elem *CurveElement) freeze() {
	if elem.isInf() {
		return // already frozen by def
	}
	elem.DataX.freeze()
	elem.DataY.freeze()
	return
}

func (elem *CurveElement) frozen() bool {
	if (elem.isInf()) {
		return true
	}
	return elem.DataX.frozen && elem.DataY.frozen
}

// TODO: for reasons I DO NOT understand, multiplication by a scalar on a curve is pow ...?
// ALSO TODO: not sure if MulScalar ends up part of Element ...?
func (elem *CurveElement) MulScalar(n *big.Int) *CurveElement {
	result := powWindow(elem, n).(*CurveElement)
	result.freeze()
	return result
}

func (elem *CurveElement) PowZn(elemIn Element) *CurveElement {
	zrElem := elemIn.(*ZrElement)
	result := powWindow(elem, zrElem.GetInt()).(*CurveElement)
	result.freeze()
	return result
}

func (elem *CurveElement) isValid() bool {

	if elem.isInf() {
		return true
	}

	calcY2 := elem.ElemField.calcYSquared(elem.DataX)
	calcY2Check := elem.DataY.square(elem.ElemField.getTargetField().FieldOrder)

	return calcY2.IsEqual(calcY2Check)
}

func (elem *CurveElement) isEqual(cmpElem *CurveElement) bool {
	if !elem.DataX.IsEqual(cmpElem.DataX) {
		return false
	}
	return elem.DataY.IsEqual(cmpElem.DataY)
}

func (elem *CurveElement) Copy() Element {
	theCopy := elem.dup()
	theCopy.freeze()
	return theCopy
}

func (elem *CurveElement) dup() *CurveElement {
	newElem := new(CurveElement)
	newElem.ElemField = elem.ElemField
	newElem.DataX = elem.DataX.copy()
	newElem.DataY = elem.DataY.copy()
	return newElem
}

func (elem *CurveElement) SetToOne() Element {
	return &CurveElement{elem.ElemField, nil, nil}
}

func (elem *CurveElement) Mul(elemIn Element) Element {
	res := elem.mul(elemIn.(*CurveElement))
	return res
}

func (elem *CurveElement) set(in *CurveElement) {
	elem.DataX = in.DataX
	elem.DataY = in.DataY
}

func (elem *CurveElement) twiceInternal() *CurveElement {

	if !elem.frozen() {
		panic("elem input must be frozen")
	}

	// We have P1 = P2 so the tangent line T at P1 ha slope
	// lambda = (3x^2 + a) / 2y
	targetOrder := elem.ElemField.getTargetField().FieldOrder
	lambdaNumer := elem.DataX.square(targetOrder).mul(BI_THREE, targetOrder).add(elem.ElemField.a.Data, targetOrder)
	lambdaDenom := elem.DataY.add(elem.DataY, targetOrder).invert(targetOrder)
	lambda := lambdaNumer.mul(lambdaDenom, targetOrder)
	lambda.freeze()

	// x3 = lambda^2 - 2x
	x3 := lambda.square(targetOrder).sub(elem.DataX.add(elem.DataX, targetOrder), targetOrder)

	// y3 = (x - x3) lambda - y
	y3 := elem.DataX.sub(x3, targetOrder).mul(lambda, targetOrder).sub(elem.DataY, targetOrder)

	x3.freeze()
	y3.freeze()
	return &CurveElement{ elem.ElemField, x3, y3}
}

func (elem *CurveElement) mul(elemIn *CurveElement) *CurveElement {

	if !elemIn.frozen() {
		panic("elemIn param must be frozen")
	}

	if elem.isInf() {
		return elemIn
	}

	if elemIn.isInf() {
		return elem
	}

	if elem.DataX.IsEqual(elemIn.DataX) {
		if elem.DataY.IsEqual(elemIn.DataY) {
			if elem.DataY.IsEqual(BI_ZERO) {
				return &CurveElement{elem.ElemField, nil, nil}
			} else {
				return elem.twiceInternal()
			}
		}
		return &CurveElement{elem.ElemField, nil, nil}
	}

	// P1 != P2, so the slope of the line L through P1 and P2 is
	// lambda = (y2-y1)/(x2-x1)
	targetOrder := elem.ElemField.getTargetField().FieldOrder
	lambdaNumer := elemIn.DataY.sub(elem.DataY, targetOrder)
	lambdaDenom := elemIn.DataX.sub(elem.DataX, targetOrder)
	lambda := lambdaNumer.mul(lambdaDenom.invert(targetOrder), targetOrder)
	lambda.freeze()

	// x3 = lambda^2 - x1 - x2
	x3 := lambda.square(targetOrder).sub(elem.DataX, targetOrder).sub(elemIn.DataX, targetOrder)

	// y3 = (x1-x3) lambda - y1
	y3 := elem.DataX.sub(x3, targetOrder).mul(lambda, targetOrder).sub(elem.DataY, targetOrder)

	x3.freeze()
	y3.freeze()
	return &CurveElement{elem.ElemField, x3, y3}
}
