package pairing

import (
	"math/big"
	"gobdc/field"
)

type GTFiniteField struct {
	field.BaseField
	targetField *field.D2ExtensionQuadField
	pairing Mapping
}

//   public GTFiniteField(SecureRandom random, BigInteger order, PairingMap pairing, F targetField) {
func MakeGTFiniteField(order *big.Int, inPairing Mapping, targetField *field.D2ExtensionQuadField) *GTFiniteField {

	field := new(GTFiniteField)

	field.targetField = targetField
	field.FieldOrder = order
	field.LengthInBytes = field.targetField.LengthInBytes * 2
	field.pairing = inPairing

	return field
}

func (field *GTFiniteField) MakeElement(inValue field.PointElement, inPairing Mapping) *GTFiniteElement {
	return &GTFiniteElement{field, inValue, inPairing}
}

var _ field.PointElement = (*GTFiniteElement)(nil)

/*
. GTFiniteElement doesn't seem to do anything besides wrap an Element value and point back to it's pairing
. The methods setToRandom() and setFromHash(...) set the value and then call the finalPow() method of the pairing
.. not sure what that is all about ...?
 */

type GTFiniteElement struct {
	ElemField *GTFiniteField
	field.PointElement
	pairing Mapping
}


