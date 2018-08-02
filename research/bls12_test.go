package research

import (
	"testing"
	"math/big"
	"math/bits"
	"fmt"
	"github.com/paulgoleary/gopbc/field"
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
	testD2SquareCompat(t)
	testD2MulCompat(t)
	testD6MulCompat(t)
}

func testD2MulCompat(t *testing.T) {

	zfield := getCompatZField()
	fieldOrder := zfield.FieldOrder

	a0 := []big.Word {2150848578656294305, 946192937311469246, 1454339966219856702, 4435310801260134845, 12289653446695321175, 123466237115473552}
	a1 := []big.Word {16351465897918308740, 17048360080029272805, 2922859199584392595, 2313888367831359194, 11821668852635520647, 1030670309910343602}

	b0 := []big.Word {2757976678761980522, 10002086052718260986, 1089158891022391851, 16355194147715895330, 90348546901268575, 12599341247547557}
	b1 := []big.Word {3886979113136472912, 14153888031467213185, 17791937096537412481, 11802835342090144038, 17849589728680247787, 1744421994245339626}

	c0 := []big.Word {3005590030555700651, 5827967474288041150, 7305580324595144075, 13246468748014481367, 13249945985045276657, 533197653467470250}
	c1 := []big.Word {11064197072050547300, 9296904173538178738, 16867030276399080008, 5461314364182372245, 17130009961149559756, 1632655424529565098}

	qfield := field.MakeD2ExtensionQuadField(zfield)

	a := qfield.MakeElement(field.MakeModIntWords(a0, true, fieldOrder),
		field.MakeModIntWords(a1, true, fieldOrder))

	b := qfield.MakeElement(field.MakeModIntWords(b0, true, fieldOrder),
		field.MakeModIntWords(b1, true, fieldOrder))

	c := qfield.MakeElement(field.MakeModIntWords(c0, true, fieldOrder),
		field.MakeModIntWords(c1, true, fieldOrder))

	d := a.MulPoint(b)
	if !c.IsValEqual(d) {
		t.Errorf("Failed to calc correct (compatible?) d2 point multiplication: want %s, got %s", c.String(), d.String())
	}
}

func testD2SquareCompat(t *testing.T) {

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

	testMulSquare := testElem.MulPoint(testElem)
	if !testMulSquare.IsValEqual(resElem) {
		t.Errorf("Failed to calc correct d2 field square by multiplication operation: expect %s, got %s", resElem.String(), testSquare.String())
	}
}

func makeTestD6Element(d2field *field.D2ExtensionQuadField, d6field *field.D6ExtensionQuadField, fieldOrder *big.Int, elemBits [][]big.Word) *field.D6ExtensionQuadElement {

	elem0 := d2field.MakeElement(
		field.MakeModIntWords(elemBits[0], true, fieldOrder),
		field.MakeModIntWords(elemBits[1], true, fieldOrder))

	elem1 := d2field.MakeElement(
		field.MakeModIntWords(elemBits[2], true, fieldOrder),
		field.MakeModIntWords(elemBits[3], true, fieldOrder))

	elem2 := d2field.MakeElement(
		field.MakeModIntWords(elemBits[4], true, fieldOrder),
		field.MakeModIntWords(elemBits[5], true, fieldOrder))

	return d6field.MakeElement(elem0, elem1, elem2)
}

func reverse(numbers []big.Word) []big.Word {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
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
	// UGH - put the data in backwards - too lazy to change it :/
	for _, value := range aBits {
		reverse(value)
	}
	aElem := makeTestD6Element(d2field, d6field, fieldOrder, aBits)

	bBits := [][]big.Word {
		{0x002CC308BB76C8A5, 0x0140FB865405D85F, 0xE2F9526C7F275022, 0x0F1D783DA3591A2B, 0x8ACE8C45892932FA, 0x26464D2D4C6B626A},
		{0x18356E8FB95F39EA, 0xF7B67B497C7281EB, 0xA3CC1835A3B58926, 0xF6E9A88609351381, 0xC46CB62C92419181, 0x35F152E8C7DCC350},
		{0x0FC5C8320F8DB723, 0x2129A0EA3CC35128, 0x69820D01A48D747C, 0x8773D5E308A31B00, 0xF1276EECE212270B, 0xD274448C84EA04BE},
		{0x0A09F11F6EB75700, 0x12274FEE19416F45, 0x71A8D0F125FC1639, 0x52BB95C2114506F5, 0x2BDBE512EC37C426, 0x9713AC654B2A72B0},
		{0x18964B178D748917, 0xE1A44FF35BD33D0C, 0x128AFEC8C2865D36, 0x7B02625FD1080E50, 0xB9F4F169E865E1BB, 0x0A40280754D41169},
		{0x0A7C1D8825BDBFB8, 0x486B3A1B4DD4F17E, 0xF89083E47E84B76A, 0x57FF533E6BABC3C7, 0x2A304DF533AAF949, 0x90D852EF7BDC8DCD}}
	for _, value := range bBits {
		reverse(value)
	}
	bElem := makeTestD6Element(d2field, d6field, fieldOrder, bBits)

	cBits := [][]big.Word {
		{0x0E02E11185BB5996, 0xE2E83F1877D6C026, 0x2BC1245FA12E0CA1, 0x6790093C24D543E1, 0x41531130F7886B68, 0x644ACE657C0A0857},
		{0x01F47ADBEA84C7D0, 0xC702A24E0164D90A, 0xF804BECC64CB1AFA, 0x6E4CC629A26E2D5F, 0x92CF8FB8C0B000FD, 0x616BE2F9804FD6C9},
		{0x129FD7B19E42DDDA, 0x2C8CF59BE2E98682, 0xF2F5AC7D93875B50, 0x1ECEC143A45225BE, 0x3DD96BF0F91F53FF, 0x459163BF524F3BEB},
		{0x0801117D8C371941, 0x3167AD4AAFE2008B, 0xCB792E686EB4E527, 0x953037874025F38A, 0xAC419A7BB8EF2FC9, 0x7AD7EC6E5BCBDECE},
		{0x098669668910CF53, 0xD99926DEE9A42202, 0x025F8E320063AC85, 0x46126F5FDB67341B, 0xFFBCB8E544915EBA, 0x6E6D613691555885},
		{0x069F1C5B3FEAA568, 0x1821C18CAB196B7E, 0x54F4867CFD0D6472, 0xC8F1D747852AB8D7, 0x60B5722198C2EE4E, 0xA29B971110C4D271}}
	for _, value := range cBits {
		reverse(value)
	}
	cElem := makeTestD6Element(d2field, d6field, fieldOrder, cBits)

	dElem := aElem.MulPoint(bElem)

	if !cElem.IsValEqual(dElem) {
		t.Errorf("Incorrect d6 element multiplication: want %s, got %s", cElem.String(), dElem.String())
	}
}