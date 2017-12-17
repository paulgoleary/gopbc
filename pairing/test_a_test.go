package pairing

import (
	"fmt"
	"gobdc/field"
	"math/big"
	"testing"
)

func testPoint(t *testing.T, thePoint field.PointLike, strExpectX string, strExpectY string) {

	checkCoord := func(cmp *field.BigInt, expString string) {
		expInt := field.MakeBigIntStr(expString, false)
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

	// these two curve point exponentiation checks are equivalent to public key derivation
	biTest1 := big.Int{}
	biTest1.SetString("276146606970621369032156664792541580771690346936", 10)
	elemTest1 := pairing.Zr.NewElement(&biTest1)
	powTest1 := pairing.G1.GetGen().PowZn(elemTest1)
	testPoint(t, powTest1,
		"2280014378744220144146373205831932526719685024545487661471834655738123196933971699437542834115250416780965121862860444719075976277314039181516434962834201",
		"5095219617050150661236292739445977344231341334112418835906977843435788249723740037212827151314561006651269984991436149205169409973600265455370653168534480")

	biTest2 := big.Int{}
	biTest2.SetString("723726390432394136624550768129737051661740488013", 10)
	elemTest2 := pairing.Zr.NewElement(&biTest2)
	powTest2 := pairing.G1.GetGen().PowZn(elemTest2)
	testPoint(t, powTest2,
		"3179305015224135534600913697529322474159056835318317733023669075672623777446135077445204913153064372784513253897383442505556687490792970594174506652914922",
		"6224780198151226873973489249032341465791104108959675858554195681300248102506693704394812888080668305563719500114924024081270783773410808141172779117133345")

	pairing.mapping.pairing(pairing.G1.GetGen(), pairing.G1.GetGen())

	println(fmt.Sprintf("Successfully made type A pairing: %T", pairing))
}
