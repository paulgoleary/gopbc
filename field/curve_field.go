package field

import (
	"math/big"
	"log"
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

	dataX := new(ModInt)
	dataX.setBytes(xBytes)
	dataX.m = field.getTargetField().FieldOrder

	dataY := new(ModInt)
	dataY.setBytes(yBytes)
	dataY.m = field.getTargetField().FieldOrder

	elem := &CurveElement{ &field.CurveParams, PointLike{dataX, dataY}}

	// needs to be frozen before validation
	elem.freeze()
	if !elem.isValid() {
		elem.setInf()
	}
	return elem
}

// general curve is y^2 = x^3 + ax + b
func (params *CurveParams) calcYSquared(xIn *ModInt) *ModInt {
	if !xIn.frozen {
		panic("xIn needs to be frozen")
	}
	validateModulo(params.getTargetField().FieldOrder, xIn.m)
	return xIn.Square().Add(params.a.Data).Mul(xIn).Add(params.b.Data)
}

// TODO: needs to account for sign
func (field *CurveField) newElementFromX(x *big.Int) *CurveElement {

	copyX := CopyFrom(x, true, field.getTargetField().FieldOrder)
	calcY2 := field.calcYSquared(copyX)
	dataY := calcY2.sqrt()

	elem := CurveElement{&field.CurveParams, PointLike{copyX, dataY}}

	elem.freeze()
	return &elem
}

func (field *CurveField) newElementFromStrings(xStr string, yStr string) *CurveElement {
	targetOrder := field.getTargetField().FieldOrder
	return &CurveElement{&field.CurveParams,
	PointLike{MakeModIntStr(xStr, true, targetOrder), MakeModIntStr(yStr, true, targetOrder)}}
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

var _ PointElement = (*CurveElement)(nil)
var _ PowElement = (*CurveElement)(nil)

func (elem *CurveElement) getTargetOrder() *big.Int {
	return elem.elemParams.getTargetField().FieldOrder
}

func (elem *CurveElement) NegateY() PointElement {
	if elem.isInf() {
		return &CurveElement{elem.elemParams, PointLike{nil, nil}}
	}
	elem.PointLike.freeze() // make sure we're frozen
	yNeg := elem.dataY.Negate()
	return &CurveElement{elem.elemParams, PointLike{elem.dataX, yNeg}}
}

func (elem *CurveElement) Invert() PointElement {
	return nil // TODO!
}

func (elem *CurveElement) Square() PointElement {
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

func validateModulo( mod1 *big.Int, mod2 *big.Int) {
	// TODO: this is intentionally pointer comparison because we expect the ModInt m's to point to the same object
	// need to think about this tho ...
	if mod1 == nil || mod1 != mod2 {
		log.Panicf("Field components must have valid and equal modulo")
	}
}

func (elem *CurveElement) isValid() bool {

	if elem.isInf() {
		return true
	}

	validateModulo(elem.dataX.m, elem.dataY.m)

	calcY2 := elem.elemParams.calcYSquared(elem.dataX)
	calcY2Check := elem.dataY.Square()

	return calcY2.IsValEqual(calcY2Check)
}

func (elem *CurveElement) isEqual(cmpElem *CurveElement) bool {
	if !elem.dataX.IsValEqual(cmpElem.dataX) {
		return false
	}
	return elem.dataY.IsValEqual(cmpElem.dataY)
}

func (elem *CurveElement) CopyPow() PowElement {
	theCopy := elem.dup()
	theCopy.freeze()
	return theCopy
}

func (elem *CurveElement) dup() *CurveElement {
	newElem := new(CurveElement)
	newElem.elemParams = elem.elemParams
	newElem.dataX = elem.dataX.Copy()
	newElem.dataY = elem.dataY.Copy()
	return newElem
}

func (elem *CurveElement) MakeOnePow() PowElement {
	return &CurveElement{elem.elemParams, PointLike{nil, nil}}
}

func (elem *CurveElement) MulPoint(elemIn PointElement) PointElement {
	res := elem.mul(elemIn.(*CurveElement))
	return res
}

func (elem *CurveElement) MulPow(elemIn PowElement) PowElement {
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
	lambdaNumer := elem.dataX.Square().Mul(MI_THREE).Add(elem.elemParams.a.Data)
	lambdaDenom := elem.dataY.Add(elem.dataY).Invert()
	lambda := lambdaNumer.Mul(lambdaDenom)
	lambda.Freeze()

	// x3 = lambda^2 - 2x
	x3 := lambda.Square().Sub(elem.dataX.Add(elem.dataX))

	// y3 = (x - x3) lambda - y
	y3 := elem.dataX.Sub(x3).Mul(lambda).Sub(elem.dataY)

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

	if elem.dataX.IsValEqual(elemIn.dataX) {
		if elem.dataY.IsValEqual(elemIn.dataY) {
			if elem.dataY.IsValEqual(MI_ZERO) {
				return &CurveElement{elem.elemParams, PointLike{nil, nil}}
			} else {
				return elem.twiceInternal()
			}
		}
		return &CurveElement{elem.elemParams, PointLike{nil, nil}}
	}

	// P1 != P2, so the slope of the line L through P1 and P2 is
	// lambda = (y2-y1)/(x2-x1)
	lambdaNumer := elemIn.dataY.Sub(elem.dataY)
	lambdaDenom := elemIn.dataX.Sub(elem.dataX)
	lambda := lambdaNumer.Mul(lambdaDenom.Invert())
	lambda.Freeze()

	// x3 = lambda^2 - x1 - x2
	x3 := lambda.Square().Sub(elem.dataX).Sub(elemIn.dataX)

	// y3 = (x1-x3) lambda - y1
	y3 := elem.dataX.Sub(x3).Mul(lambda).Sub(elem.dataY)

	x3.Freeze()
	y3.Freeze()
	return &CurveElement{elem.elemParams, PointLike {x3, y3}}
}
