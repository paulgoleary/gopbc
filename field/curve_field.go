package field

import "math/big"

type CurveField struct {
	a          *ZrElement
	b          *ZrElement
	order      *big.Int
	cofactor   *big.Int
	gen        *CurveElement // TODO: not sure here...
	genNoCofac *CurveElement
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
	field.gen = field.genNoCofac.dup().MulScalar(field.cofactor)
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

	/*
		//if point does not lie on curve, set it to O
		if (!isValid())
			setToZero();

		return len;
	*/
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

// CurveElement

// TODO: Make function?

// validate that Curve Element satisfies Element
var _ Element = (*CurveElement)(nil)

func (elem *CurveElement) IsInf() bool {
	return elem.DataY == nil && elem.DataY == nil
}

func (elem *CurveElement) SetInf() {
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

// TODO: for reasons I DO NOT understand, multiplication by a scalar on a curve is pow ...?
// ALSO TODO: not sure if MulScalar ends up part of Element ...?
func (elem *CurveElement) MulScalar(n *big.Int) *CurveElement {
	return powWindow(elem, n).(*CurveElement)
}

/*
	Element interface:
		Copy() Element
		Mul(*Element) Element
		SetToOne() Element
*/

func (elem *CurveElement) dup() *CurveElement {
	newElem := new(CurveElement)
	newElem.ElemField = elem.ElemField
	newElem.DataX = elem.DataX.copy()
	newElem.DataY = elem.DataY.copy()
	return newElem
}

func (elem *CurveElement) Copy() Element {
	return elem.dup()
}

func (elem *CurveElement) SetToOne() Element {
	elem.SetInf()
	return elem
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
	// We have P1 = P2 so the tangent line T at P1 ha slope
	// lambda = (3x^2 + a) / 2y
	targetOrder := elem.ElemField.getTargetField().FieldOrder
	lambdaNumer := elem.DataX.copy().square(targetOrder).mul(BI_THREE, targetOrder).add(elem.ElemField.a.Data, targetOrder)
	lambdaDenom := elem.DataY.copy().add(elem.DataY, targetOrder).invert(targetOrder)
	lambda := lambdaNumer.mul(lambdaDenom, targetOrder)

	// x3 = lambda^2 - 2x
	x3 := lambda.copy().square(targetOrder).sub(elem.DataX.copy().add(elem.DataX, targetOrder), targetOrder)

	// y3 = (x - x3) lambda - y
	y3 := elem.DataX.copy().sub(x3, targetOrder).mul(lambda, targetOrder).sub(elem.DataY, targetOrder)

	elem.DataX = x3
	elem.DataY = y3
	return elem
}

func (elem *CurveElement) mul(elemIn *CurveElement) *CurveElement {

	if elem.IsInf() {
		elem.set(elemIn)
		return elem
	}

	if elemIn.IsInf() {
		return elem
	}

	if elem.DataX.IsEqual(elemIn.DataX) {
		if elem.DataY.IsEqual(elemIn.DataY) {
			if elem.DataY.IsEqual(BI_ZERO) {
				elem.SetInf()
				return elem
			} else {
				elem.twiceInternal()
				return elem
			}
		}
		elem.SetInf()
		return elem
	}

	// P1 != P2, so the slope of the line L through P1 and P2 is
	// lambda = (y2-y1)/(x2-x1)
	targetOrder := elem.ElemField.getTargetField().FieldOrder
	lambdaNumer := elemIn.DataY.copy().sub(elem.DataY, targetOrder)
	lambdaDenom := elemIn.DataX.copy().sub(elem.DataX, targetOrder)
	lambda := lambdaNumer.mul(lambdaDenom.invert(targetOrder), targetOrder)

	// x3 = lambda^2 - x1 - x2
	x3 := lambda.copy().square(targetOrder).sub(elem.DataX, targetOrder).sub(elemIn.DataX, targetOrder)

	// y3 = (x1-x3) lambda - y1
	y3 := elem.DataX.copy().sub(x3, targetOrder).mul(lambda, targetOrder).sub(elem.DataY, targetOrder)

	elem.DataX = x3
	elem.DataY = y3

	return elem
}
