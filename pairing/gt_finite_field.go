package pairing

import (
	"math/big"
	"gopbc/field"
)

type GTFiniteField struct {
	field.BaseField
	targetField *field.D2ExtensionQuadField
	pairing Mapping
}

//   public GTFiniteField(SecureRandom random, BigInteger order, PairingMap Pairing, F targetField) {
func MakeGTFiniteField(order *big.Int, inPairing Mapping, targetField *field.D2ExtensionQuadField) *GTFiniteField {

	field := new(GTFiniteField)

	field.targetField = targetField
	field.FieldOrder = order
	field.LengthInBytes = field.targetField.LengthInBytes
	field.pairing = inPairing

	return field
}

func (field *GTFiniteField) MakeElement(inValue field.PointElement, inPairing Mapping) *GTFiniteElement {
	return &GTFiniteElement{field, inValue, inPairing}
}

func (field *GTFiniteField) MakeElementFromBytes(elemBytes []byte, inPairing Mapping) *GTFiniteElement {
	elem := field.targetField.MakeElementFromBytes(elemBytes)
	return &GTFiniteElement{field, elem, inPairing}
}

var _ field.PointElement = (*GTFiniteElement)(nil)

/*
. GTFiniteElement doesn't seem to do anything besides wrap an Element value and point back to it's Pairing
. The methods setToRandom() and setFromHash(...) set the value and then call the finalPow() method of the Pairing
.. not sure what that is all about ...?
 */

type GTFiniteElement struct {
	ElemField *GTFiniteField
	field.PointElement
	pairing Mapping
}


