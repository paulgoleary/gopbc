package research

import (
	"testing"
	"github.com/paulgoleary/gopbc/field"
	"math/big"
	)

func makeTestFields() (*field.CurveField, *field.D2ExtensionQuadField) {
	Fq := field.MakeZField(big.NewInt(59))

	curveField := field.MakeCurveField(
		Fq.NewOneElement(),
		Fq.NewZeroElement(),
		big.NewInt(5),
		nil,
		nil )

	extField := field.MakeD2ExtensionQuadField(Fq)

	return curveField, extField
}
func TestToyField(t *testing.T) {

	curveField, _ := makeTestFields()

	// for testing purposes I want the *other* point at x = 25
	testGen := curveField.MakeElementFromX(big.NewInt(25)).Invert()
	field.Trace(testGen)

	testOrder := field.PointElement(testGen)
	for i := 0; i < 4; i++ {
		testOrder = testOrder.Add(testGen)
		field.Trace(testOrder)
	}
	if !testOrder.(*field.CurveElement).IsInf() {
		t.Errorf("Should have cycled back to O element of subgroup")
	}
}

func TestToyWeilPairing(t *testing.T) {

	curveField, _ := makeTestFields()

	testP := curveField.MakeElement(big.NewInt(25), big.NewInt(30))

	testQ := curveField.MakeElement( big.NewInt(-25), big.NewInt(30) )
	field.Trace(testQ)

	R := curveField.MakeElement(big.NewInt(40), big.NewInt(54))

	S0 := curveField.MakeElement(big.NewInt(48), big.NewInt(55))
	field.Trace(S0)
	S1 := curveField.MakeElement(big.NewInt(28), big.NewInt(51))
	field.Trace(S1)

	PR := testP.Add(R)
	field.Trace(PR)

	testN := field.MakeModInt(3339376, true, big.NewInt(59))
	testD := field.MakeModInt(-3600, true, big.NewInt(59)).Invert()
	testX := testN.Mul(testD)
	field.Trace(testX)
}
