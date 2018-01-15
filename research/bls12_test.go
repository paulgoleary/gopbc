package research

import (
	"testing"
	"math/big"
	"math/bits"
	"fmt"
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

func TestBLS12ParamGen(t *testing.T) {

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
