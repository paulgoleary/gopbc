package research

import (
	"testing"
	"math/big"
	"math/bits"
	"fmt"
	"gopbc/field"
)

func hammingWeight( x *big.Int ) int {

	// TODO: pre-init elsewhere ...?
	m := make(map[uint8]int)
	for i := 0; i < 256; i += 1 {
		ui := uint8(i)
		m[ui] = bits.OnesCount8(ui)
	}

	cntBits := 0
	xBytes := x.Bytes()
	for _, xx := range xBytes {
		cntBits += m[uint8(xx)]
	}
	return cntBits
}

var BI_ZERO = big.NewInt(0)
var BI_ONE = big.NewInt(1)
var BI_TWO = big.NewInt(2)
var BI_THREE = big.NewInt(3)
var BI_FOUR = big.NewInt(4)
var BI_N_TWO = big.NewInt(-2)
var BI_N_ONE = big.NewInt(-1)
var BI_FIVE = big.NewInt(5)
var BI_SIX = big.NewInt(6)

type polyTerm struct {
	exp *big.Int
	coef *big.Int
}

type polyDef []polyTerm

func (def polyDef) eval(x *big.Int) *big.Int {
	res := big.NewInt(0)
	for _, term := range def {
		xx := big.Int{}
		if term.exp != nil {
			xx.Exp(x, term.exp, nil)
		}
		if term.coef != nil {
			xx.Mul(&xx, term.coef)
		}
		res.Add(res, &xx)
	}
	return res
}

// p = (x^2 - 2x + 1) * (x^4 - x^2 + 1)/3 + x
func relicCalc(t0 *big.Int) *big.Int {
	// bn_sqr(t1, t0);
	t1 := new(big.Int).Set(t0)
	t1.Mul(t1, t0) // t1 = t0 ^ 2

	// bn_sqr(p, t1);
	p := new(big.Int).Set(t1)
	p.Mul(p, t1) // p = t1 * t1 = t1 ^ 2 = t0 ^ 4

	// bn_sub(p, p, t1);
	p.Sub(p, t1) // p = p - t1 = t0^4 - t0^2

	// bn_add_dig(p, p, 1);
	p.Add(p, BI_ONE) // p = t0^4 - t0^2 + 1
	println(fmt.Sprintf("%x", p))

	// bn_sub(t1, t1, t0);
	t1.Sub(t1, t0) // t1 = t0^2 - t0
	// bn_sub(t1, t1, t0);
	t1.Sub(t1, t0) // t1 = t0^2 - 2*t0

	// bn_add_dig(t1, t1, 1);
	t1.Add(t1, BI_ONE) // t1 = t0^2 - 2*t0 + 1
	println(fmt.Sprintf("%x", t1))

	// bn_mul(p, p, t1);
	p.Mul(p, t1)

	// bn_div_dig(p, p, 3);
	p.Div(p, BI_THREE)

	// bn_add(p, p, t0);
	p.Add(p, t0)

	return p
}

func TestBLS12Params(t *testing.T) {

	testXString := "-d201000000010000"

	/* x = -(2^63 + 2^62 + 2^60 + 2^57 + 2^48 + 2^16). */
	xPoly := polyDef {
		{big.NewInt(63), nil},
		{big.NewInt(62), nil},
		{big.NewInt(60), nil},
		{big.NewInt(57), nil},
		{big.NewInt(48), nil},
		{big.NewInt(16), nil}}
	calcX := xPoly.eval(BI_TWO)
	calcX.Neg(calcX)
	calcXString := fmt.Sprintf("%x", calcX)
	if testXString != calcXString {
		t.Errorf("Calc'd wrong value for X: want %s, got %s", testXString, calcXString)
	}

	xhw := hammingWeight(calcX)
	if  xhw != 6  {
		t.Errorf("Got wrong Hamming weight for calc'd X: got %v, want 6", xhw)
	}

	// pRelic := relicCalc(calcX)
	// println(fmt.Sprintf("%x", pRelic))

	// p = (x^2 - 2x + 1) * (x^4 - x^2 + 1)/3 + x
	pPoly1 := polyDef {
		{BI_TWO, nil},
		{BI_ONE, BI_N_TWO},
		{BI_ZERO, nil} }
	pPoly2 := polyDef {
		{BI_FOUR, nil},
		{BI_TWO, BI_N_ONE},
		{BI_ZERO, nil} }

	p2 := pPoly2.eval(calcX)
	p1 := pPoly1.eval(calcX)

	p := p1.Mul(p1, p2)
	p.Div(p, BI_THREE)
	p.Add(p, calcX)
	calcPString := fmt.Sprintf("%x", p)

	testPString := "1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab"

	if calcPString != testPString {
		t.Errorf("Calc'd wrong value for P: want %s, got %s", testPString, calcPString)
	}
}

func getCompatZField() *field.ZField {
	fieldOrder := new(big.Int)
	fieldOrder.SetString("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab", 16)
	return field.MakeZField(fieldOrder)
}

func TestBLS12Fields(t *testing.T) {
	testD2MulCompat(t)
	testD6MulCompat(t)
}

func testD2MulCompat(t *testing.T) {

	zfield := getCompatZField()
	fieldOrder := zfield.FieldOrder

	v0 := []big.Word {
		15980986403153949556,
		8312704975630804322,
		14168571785075134642,
		14165629176482031449,
		7389832673998068853,
		1781118126119810123	}

	v1 := []big.Word {
		17909906142013796791,
		15686367500387818584,
		16823448455851608339,
		6081957736848945043,
		9598562759131955863,
		991220194111438160 }

	p0 := field.MakeModIntWords(v0, true, fieldOrder)
	p1 := field.MakeModIntWords(v1, true, fieldOrder)

	qfield := field.MakeD2ExtensionQuadField(zfield)

	testElem := qfield.MakeElement(p0, p1)

	r0 := []big.Word {
		14933914844065159370,
		2122104236263514281,
		4151227062241471315,
		15249847645077438991,
		10892634782589633608,
		271899024200359973 }

	r1 := []big.Word {
		13878994217214973151,
		10552726046066292122,
		14454611607008347041,
		15499325862293857853,
		10393135513138888064,
		1584683573066562359	}

	resElem := qfield.MakeElement(
		field.MakeModIntWords(r0, true, fieldOrder),
		field.MakeModIntWords(r1, true, fieldOrder))

	testSquare := testElem.Square()
	if !testSquare.IsValEqual(resElem) {
		t.Errorf("Failed to calc correct d2 field square operation: expect %s, got %s", resElem.String(), testSquare.String())
	}
}

func testD6MulCompat(t *testing.T) {

	zfield := getCompatZField()
	fieldOrder := zfield.FieldOrder

	d2field := field.MakeD2ExtensionQuadField(zfield)
	d6field := field.MakeD6ExtensionQuadField(d2field)

	aBits := [][]big.Word {
		{0x01B6A3E375059E90, 0xAA8D9EB246EDFE57, 0x3D8D63ACA99041BD, 0x142EDA7E6D302B3E, 0x0D218D769D6EFEBE, 0x1DD95962F87655A1},
		{0x0E4DAD2FA6BF6FB2, 0xA40F013005A3CE87, 0x201C9536067F1EDA, 0x289014F4BFB5D593, 0xEC97F12AAD184AE5, 0xE2EC1399B255A184},
		{0x06094E21296EBAD8, 0x5D0C6207C82C9173, 0x86EACA537FE44AEA, 0xA044254B549835BF, 0x2817A10CAE069A53, 0x006B13E4173FAFC2},
		{0x051AFDD5F166623A, 0x6B16363F60528512, 0xFD3645D67F09E01A, 0x7157A501760DCDF0, 0x5D0F5CBF3AB0156E, 0xF242B1DEE7592FC6},
		{0x17C6BA91C4B272E6, 0xE37C6693355B0197, 0x8C8FA99B12B47AE7, 0xC84C0AEB02DFBD5A, 0xC5379FAE8D1FBAD6, 0x7C3EDFB568979DCA},
		{0x18EC8049302B36A0, 0x94C36E460CE3FED6, 0x755A1968D3A397ED, 0x90278397D7124D8A, 0x9F9D23F1CEB0AA99, 0xB1AD412B43192E1D}}

	aElem0 := d2field.MakeElement(
		field.MakeModIntWords(aBits[0], true, fieldOrder),
		field.MakeModIntWords(aBits[1], true, fieldOrder))

	aElem1 := d2field.MakeElement(
		field.MakeModIntWords(aBits[2], true, fieldOrder),
		field.MakeModIntWords(aBits[3], true, fieldOrder))

	aElem2 := d2field.MakeElement(
		field.MakeModIntWords(aBits[4], true, fieldOrder),
		field.MakeModIntWords(aBits[5], true, fieldOrder))

	aElem := d6field.MakeElement(aElem0, aElem1, aElem2)
	field.Trace(aElem)

}