package field

import "math/big"

type QuadraticField struct {
	BaseField
	targetField *ZrField
}

// D2ExtensionQuadElement

var _ Element = (*D2ExtensionQuadElement)(nil)
var _ PointElement = (*D2ExtensionQuadElement)(nil)

type D2ExtensionQuadElement struct {
	ElemField *D2ExtensionQuadField
	PointLike
}

func (elem *D2ExtensionQuadElement) X() *BigInt {
	return elem.DataX
}

func (elem *D2ExtensionQuadElement) Y() *BigInt {
	return elem.DataY
}

// TODO !!!
func (elem *D2ExtensionQuadElement) Negate() PointElement {
	return elem
}

func (elem *D2ExtensionQuadElement) Copy() Element {
	theCopy := elem.dup()
	theCopy.freeze()
	return theCopy
}

func (elem *D2ExtensionQuadElement) dup() *D2ExtensionQuadElement {
	newElem := new(D2ExtensionQuadElement)
	newElem.ElemField = elem.ElemField
	newElem.DataX = elem.DataX.copy()
	newElem.DataY = elem.DataY.copy()
	return newElem
}

func (elem *D2ExtensionQuadElement) SetToOne() Element {
	return &D2ExtensionQuadElement{elem.ElemField, PointLike{BI_ONE, BI_ZERO}}
}

func (elem *D2ExtensionQuadElement) Square() Element {
	/*
	Element e0 = this.x.duplicate();
    Element e1 = this.x.duplicate();
    e0.add(this.y).mul(e1.sub(this.y));
    e1.set(this.x).mul(this.y).twice();
    this.x.set(e0);
    this.y.set(e1);
    return this;
	 */

	targetOrder := elem.ElemField.targetField.FieldOrder // TODO: verify
	e0 := elem.DataX.Add(elem.DataY, targetOrder).Mul(elem.DataX.Sub(elem.DataY, targetOrder), targetOrder)
	e1 := elem.DataX.Mul(elem.DataY, targetOrder).Mul(BI_TWO, targetOrder)
	return elem.ElemField.MakeElement(e0, e1)
}

func (elem *D2ExtensionQuadElement) Mul(elemIn Element) Element {
	/*
	DegreeTwoExtensionQuadraticElement element = (DegreeTwoExtensionQuadraticElement)e;

	Element e0 = this.x.duplicate();
	Element e1 = element.x.duplicate();
    Element e2 = this.x.getField().newElement();

    e2.set(e0.add(this.y)).mul(e1.add(element.y));
    e0.set(this.x).mul(element.x);
    e1.set(this.y).mul(element.y);
    e2.sub(e0);

    this.x.set(e0).sub(e1);
    this.y.set(e2).sub(e1);

    return this;
	 */

	targetOrder := elem.ElemField.targetField.FieldOrder // TODO - verify !
	d2xqeIn := elemIn.(*D2ExtensionQuadElement) // curses!!! was hoping to avoid this :/
	e2 := elem.DataX.Add(elem.DataY, targetOrder).Mul(d2xqeIn.DataX.Add(d2xqeIn.DataY, targetOrder), targetOrder)
	e0 := elem.DataX.Mul(d2xqeIn.DataX, targetOrder)
	e1 := elem.DataY.Mul(d2xqeIn.DataY, targetOrder)
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