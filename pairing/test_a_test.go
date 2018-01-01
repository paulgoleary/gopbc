package pairing

import (
	"gobdc/field"
	"math/big"
	"testing"
)

func testPoint(t *testing.T, thePoint field.PointElement, strExpectX string, strExpectY string) {

	checkCoord := func(cmp *field.ModInt, expString string) {
		expInt := field.MakeModIntStr(expString, false, nil)
		if !expInt.IsValEqual(cmp) {
			t.Errorf("Wrong value for point coord comparision, got: %s, want: %s.", cmp.String(), expString)
		}
	}
	checkCoord(thePoint.X(), strExpectX)
	checkCoord(thePoint.Y(), strExpectY)
}

// func MakeTypeAPairing(params *PairingParameters) *TypeAPairing {
func TestMakeTypeAPairing(t *testing.T) {
	pairingParms := getCompatParams()
	typeAPairing := MakeTypeAPairing(pairingParms)

	// test compatibility with PBC
	// TODO: move to a more explicit compat test?
	testPoint(t, typeAPairing.G1.GetGen(),
		"7852334875614213225969535005319230321249629225894318783946607976937179571030765324627135523985138174020408497250901949150717492683934959664497943409406486",
		"8189589736511278424487290408486860952887816120897672059241649987466710766123126805204101070682864313793496226965335026128263318306025907120292056643404206")

	// these two curve point exponentiation checks are equivalent to public key derivation
	biTest1 := big.Int{}
	biTest1.SetString("276146606970621369032156664792541580771690346936", 10)
	elemTest1 := typeAPairing.Zr.NewElement(&biTest1)
	powTest1 := typeAPairing.G1.GetGen().PowZn(elemTest1.GetValue())
	testPoint(t, powTest1,
		"2280014378744220144146373205831932526719685024545487661471834655738123196933971699437542834115250416780965121862860444719075976277314039181516434962834201",
		"5095219617050150661236292739445977344231341334112418835906977843435788249723740037212827151314561006651269984991436149205169409973600265455370653168534480")

	biTest2 := big.Int{}
	biTest2.SetString("723726390432394136624550768129737051661740488013", 10)
	elemTest2 := typeAPairing.Zr.NewElement(&biTest2)
	powTest2 := typeAPairing.G1.GetGen().PowZn(elemTest2.GetValue())
	testPoint(t, powTest2,
		"3179305015224135534600913697529322474159056835318317733023669075672623777446135077445204913153064372784513253897383442505556687490792970594174506652914922",
		"6224780198151226873973489249032341465791104108959675858554195681300248102506693704394812888080668305563719500114924024081270783773410808141172779117133345")

	Z := typeAPairing.mapping.pairing(typeAPairing.G1.GetGen(), typeAPairing.G1.GetGen())
	testPoint(t, Z,
		"8427709406215883227839951717300804898047197599691179303852509060811239041103323092291271103964136225006329684795508648068486825455284982044949366859105368",
		"4452853254007027991300077621238640815837993300635305842991137231634336596445454707932906291400921695774864920057621535330198044334793069508782439758034655")

	// e(g, g)^ab
	testPairingPow := Z.Pow(elemTest1.Mul(elemTest2.ModInt))

	// e(g^a, g^b)
	testPowPairing := typeAPairing.mapping.pairing(powTest1, powTest2)

	// e(g^a, g^b) == e(g, g)^ab - pairings FTW!
	if !testPairingPow.IsValEqual(testPowPairing) {
		t.Errorf("Expected bilinear pairings to be equivalent: %s, %s", testPowPairing.String(), testPairingPow.String())
	}
}
