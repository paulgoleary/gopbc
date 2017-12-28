package field

import (
	"math/big"
)

type CurveField struct {
	CurveParams
	cofactor   *big.Int      // TODO: do we need this ...?
	gen        *CurveElement // TODO: not sure here...
	genNoCofac *CurveElement // TODO: don't need this ...?
}

type CurveParams struct {
	BaseField
	a          *ZrElement
	b          *ZrElement
}

type CurveElement struct {
	elemParams *CurveParams
	PointLike
}

// CurveField

// TODO: JPBC (PBC?) handles case w/o bytes and cofactor
func (field *CurveField) initGenFromBytes(genNoCofacBytes *[]byte) {
	newGenNoCoFac := field.newElementFromBytes(genNoCofacBytes)
	field.genNoCofac = newGenNoCoFac
	field.gen = field.genNoCofac.MulScalar(field.cofactor)
	if !field.gen.isValid(){
		panic("Curve field generator needs to be valid")
	}
}

func (field *CurveField) GetGen() *CurveElement {
	return field.gen
}

func (curveParams *CurveParams) getTargetField() *ZrField {
	return curveParams.a.ElemField
}

func (field *CurveField) newElementFromBytes(elemBytes *[]byte) *CurveElement {

	xBytes := (*elemBytes)[:field.getTargetField().LengthInBytes]
	yBytes := (*elemBytes)[field.getTargetField().LengthInBytes:]

	dataX := new(BigInt)
	dataX.setBytes(xBytes)

	dataY := new(BigInt)
	dataY.setBytes(yBytes)

	elem := &CurveElement{ &field.CurveParams, PointLike{dataX, dataY}}

	// needs to be frozen before validation
	elem.freeze()
	if !elem.isValid() {
		elem.setInf()
	}
	return elem
}

// general curve is y^2 = x^3 + ax + b
func (params *CurveParams) calcYSquared(xIn *BigInt) *BigInt {
	order := params.getTargetField().FieldOrder
	if !xIn.frozen {
		panic("xIn needs to be frozen")
	}
	return xIn.Square(order).Add(params.a.Data, order).Mul(xIn, order).Add(params.b.Data, order)
}

// TODO: needs to account for sign
func (field *CurveField) newElementFromX(x *big.Int) *CurveElement {

	copyX := CopyFrom(x, true)
	calcY2 := field.calcYSquared(copyX)
	dataY := calcY2.sqrt(field.getTargetField().FieldOrder)

	elem := CurveElement{&field.CurveParams, PointLike{copyX, dataY}}

	elem.freeze()
	return &elem
}

func (field *CurveField) newElementFromStrings(xStr string, yStr string) *CurveElement {
	return &CurveElement{&field.CurveParams, PointLike{MakeBigIntStr(xStr, true), MakeBigIntStr(yStr, true)}}
}

func getLengthInBytes( field *CurveField ) int {
	return field.getTargetField().LengthInBytes * 2
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
	field.FieldOrder = order
	field.cofactor = cofactor
	field.LengthInBytes = getLengthInBytes(field)
	field.initGenFromBytes(genNoCofacBytes)

	return field
}

// TODO: need to reconcile this and the other make function - not sure I both ...?
// make minimal field for testing purposes - TODO: might need a generator?
func makeTestCurveField(a *big.Int, b *big.Int, r *big.Int, q *big.Int) *CurveField {

	zfield := MakeZrField(q)

	cfield := new(CurveField)
	cfield.a = zfield.NewElement(a)
	cfield.b = zfield.NewElement(b)
	cfield.FieldOrder = r
	cfield.LengthInBytes = getLengthInBytes(cfield)

	return cfield
}

// CurveElement

// TODO: Make function?

// validate that CurveElement satisfies Element
var _ PointElement = (*CurveElement)(nil)

func (elem *CurveElement) getTargetOrder() *big.Int {
	return elem.elemParams.getTargetField().FieldOrder
}

func (elem *CurveElement) NegateP() PointElement {
	if elem.isInf() {
		return &CurveElement{elem.elemParams, PointLike{nil, nil}}
	}
	elem.PointLike.freeze() // make sure we're frozen
	yNeg := elem.dataY.Negate(elem.getTargetOrder())
	return &CurveElement{elem.elemParams, PointLike{elem.dataX, yNeg}}
}

func (elem *CurveElement) InvertP() PointElement {
	return nil // TODO!
}

func (elem *CurveElement) Square() Element {
	// TODO !
	return nil
}

func (elem *CurveElement) isInf() bool {
	return elem.dataY == nil && elem.dataY == nil
}

func (elem *CurveElement) setInf() {
	elem.dataX = nil
	elem.dataY = nil
}

// don't return elem to emphasize that call mutates elem
func (elem *CurveElement) freeze() {
	if elem.isInf() {
		return // already frozen by def
	}
	elem.PointLike.freeze()
	return
}

func (elem *CurveElement) frozen() bool {
	if (elem.isInf()) {
		return true
	}
	return elem.PointLike.frozen()
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

	calcY2 := elem.elemParams.calcYSquared(elem.dataX)
	calcY2Check := elem.dataY.Square(elem.getTargetOrder())

	return calcY2.IsEqual(calcY2Check)
}

func (elem *CurveElement) isEqual(cmpElem *CurveElement) bool {
	if !elem.dataX.IsEqual(cmpElem.dataX) {
		return false
	}
	return elem.dataY.IsEqual(cmpElem.dataY)
}

func (elem *CurveElement) Copy() Element {
	theCopy := elem.dup()
	theCopy.freeze()
	return theCopy
}

func (elem *CurveElement) dup() *CurveElement {
	newElem := new(CurveElement)
	newElem.elemParams = elem.elemParams
	newElem.dataX = elem.dataX.copy()
	newElem.dataY = elem.dataY.copy()
	return newElem
}

func (elem *CurveElement) SetToOne() Element {
	return &CurveElement{elem.elemParams, PointLike{nil, nil}}
}

func (elem *CurveElement) Mul(elemIn Element) Element {
	res := elem.mul(elemIn.(*CurveElement))
	return res
}

func (elem *CurveElement) set(in *CurveElement) {
	elem.dataX = in.dataX
	elem.dataY = in.dataY
}

func (elem *CurveElement) twiceInternal() *CurveElement {

	if !elem.frozen() {
		panic("elem input must be frozen")
	}

	// We have P1 = P2 so the tangent line T at P1 ha slope
	// lambda = (3x^2 + a) / 2y
	targetOrder := elem.getTargetOrder()
	lambdaNumer := elem.dataX.Square(targetOrder).Mul(BI_THREE, targetOrder).Add(elem.elemParams.a.Data, targetOrder)
	lambdaDenom := elem.dataY.Add(elem.dataY, targetOrder).invert(targetOrder)
	lambda := lambdaNumer.Mul(lambdaDenom, targetOrder)
	lambda.Freeze()

	// x3 = lambda^2 - 2x
	x3 := lambda.Square(targetOrder).Sub(elem.dataX.Add(elem.dataX, targetOrder), targetOrder)

	// y3 = (x - x3) lambda - y
	y3 := elem.dataX.Sub(x3, targetOrder).Mul(lambda, targetOrder).Sub(elem.dataY, targetOrder)

	x3.Freeze()
	y3.Freeze()
	return &CurveElement{ elem.elemParams, PointLike {x3, y3}}
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

	if elem.dataX.IsEqual(elemIn.dataX) {
		if elem.dataY.IsEqual(elemIn.dataY) {
			if elem.dataY.IsEqual(BI_ZERO) {
				return &CurveElement{elem.elemParams, PointLike{nil, nil}}
			} else {
				return elem.twiceInternal()
			}
		}
		return &CurveElement{elem.elemParams, PointLike{nil, nil}}
	}

	// P1 != P2, so the slope of the line L through P1 and P2 is
	// lambda = (y2-y1)/(x2-x1)
	targetOrder := elem.getTargetOrder()
	lambdaNumer := elemIn.dataY.Sub(elem.dataY, targetOrder)
	lambdaDenom := elemIn.dataX.Sub(elem.dataX, targetOrder)
	lambda := lambdaNumer.Mul(lambdaDenom.invert(targetOrder), targetOrder)
	lambda.Freeze()

	// x3 = lambda^2 - x1 - x2
	x3 := lambda.Square(targetOrder).Sub(elem.dataX, targetOrder).Sub(elemIn.dataX, targetOrder)

	// y3 = (x1-x3) lambda - y1
	y3 := elem.dataX.Sub(x3, targetOrder).Mul(lambda, targetOrder).Sub(elem.dataY, targetOrder)

	x3.Freeze()
	y3.Freeze()
	return &CurveElement{elem.elemParams, PointLike {x3, y3}}
}
