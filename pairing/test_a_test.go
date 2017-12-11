package pairing

import (
	"fmt"
	"gobdc/field"
	"testing"
)

func testPoint(t *testing.T, thePoint field.PointLike, strExpectX string, strExpectY string) {

	checkCoord := func(cmp *field.BigInt, expString string) {
		expInt := field.MakeBigIntStr(expString)
		if !expInt.IsEqual(cmp) {
			t.Errorf("Wrong value for point coord comparision, got: %s, want: %s.", cmp.String(), expString)
		}
	}
	checkCoord(thePoint.X(), strExpectX)
	checkCoord(thePoint.Y(), strExpectY)
}

// func MakeTypeAPairing(params *PairingParameters) *TypeAPairing {
func TestMakeTypeAPairing(t *testing.T) {
	pairingParms := getCompatParams()
	pairing := MakeTypeAPairing(pairingParms)

	// test compatibility with PBC
	// TODO: move to a more explicit compat test?
	testPoint(t, pairing.G1.GetGen(),
		"7852334875614213225969535005319230321249629225894318783946607976937179571030765324627135523985138174020408497250901949150717492683934959664497943409406486",
		"8189589736511278424487290408486860952887816120897672059241649987466710766123126805204101070682864313793496226965335026128263318306025907120292056643404206")

	println(fmt.Sprintf("Successfully made type A pairing: %T", pairing))
}
