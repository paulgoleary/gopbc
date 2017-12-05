package pairing

import (
	"fmt"
	"gobdc/field"
	"math/big"
	"testing"
)

func getCompatParams() *PairingParameters {

	/*
	   *** this is now actually *also* cribbed from the bdc project ...
	   taken from PBC-provided param file 'a.param'
	   type a
	   q 8780710799663312522437781984754049815806883199414208211028653399266475630880222957078625179422662221423155858769582317459277713367317481324925129998224791
	   h 12016012264891146079388821366740534204802954401251311822919615131047207289359704531102844802183906537786776
	   r 730750818665451621361119245571504901405976559617
	   exp2 159
	   exp1 107
	   sign1 1
	   sign0 1
	*/

	params := PairingParameters{}
	params["type"] = "a"
	params["q"] = "8780710799663312522437781984754049815806883199414208211028653399266475630880222957078625179422662221423155858769582317459277713367317481324925129998224791"
	params["r"] = "730750818665451621361119245571504901405976559617"
	params["h"] = "12016012264891146079388821366740534204802954401251311822919615131047207289359704531102844802183906537786776"
	params["exp1"] = "107"
	params["exp2"] = "159"
	params["sign0"] = "1"
	params["sign1"] = "1"

	params["genNoCofac"] = "WLeuxaO0DxaW+oJ4vrLKgkq91prZNLGQUVoXH4gIx6AGIS7vrU7Fq3/5DfYTRHfpnOCIuo96hfRwTzUTf2+EUndlGtVaI05vjWxsIaCqKSPtq+xYpr0jaGVVwnXojhjbi0AeR/JvjiIaF9WFjSRzqEvR1WHp0LkJRrtBfNcA0k4="

	return &params
}

func testPoint(t *testing.T, thePoint field.Point, strExpectX string, strExpectY string) {

	checkCoord := func(cmp *big.Int, expString string) {
		expInt := big.Int{}
		expInt.SetString(expString, 10)
		if expInt.Cmp(cmp) != 0 {
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
