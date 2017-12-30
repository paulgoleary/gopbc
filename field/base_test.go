package field

import (
	"math/big"
	"math/rand"
	"sort"
	"testing"
)

type TestElement struct {
	data *big.Int
	mod  big.Int
}

// implement PowElement for TestElement

func (elem *TestElement) CopyPow() PowElement {
	newElem := new(TestElement)
	newElem.data = new(big.Int)
	newElem.data.SetBytes(elem.data.Bytes()) // need to make 'deep Copy' of mutable data
	newElem.mod = elem.mod
	return newElem
}

func (elem *TestElement) MulPow(mulElem PowElement) PowElement {
	in := mulElem.(*TestElement) // TODO: not a fan of this...
	ret := elem.CopyPow().(*TestElement)
	ret.data.Mul(elem.data, in.data)
	ret.data.Mod(ret.data, &elem.mod)
	return ret
}

func (elem *TestElement) MakeOnePow() PowElement {
	ret := elem.CopyPow().(*TestElement)
	ret.data.Set(ONE)
	return ret
}

func (elem *TestElement) String() string {
	return elem.data.String()
}

// validate that TestElement satisfies Element interface
var _ PowElement = (*TestElement)(nil)

func checkPowWindowBigInt(t *testing.T, testBase *TestElement, testExp *big.Int) {
	expectedBigInt := new(big.Int).Exp(testBase.data, testExp, &testBase.mod)
	checkPowWindow(t, testBase, testExp, expectedBigInt)
}

func checkPowWindow(t *testing.T, testBase *TestElement, testExp, expectedVal *big.Int) {
	var elem PowElement = testBase
	testPow := powWindow(elem, testExp)
	if testPow.(*TestElement).data.Cmp(expectedVal) != 0 {
		t.Errorf("powWindow exponent result was wrong, got: %d, want: %d.", testPow.String(), expectedVal.String())
	}
}

func makeBigInt(valStr string) (res *big.Int) {
	res = new(big.Int)
	res.SetString(valStr, 10)
	return
}

func makeRandBigInt() (res *big.Int) {
	res = new(big.Int)
	resBits := make([]byte, 20, 20) // 160 bits
	rand.Read(resBits)
	res.SetBytes(resBits)
	return
}

// TODO: maybe already implemented for big.Int?
type BigInts []*big.Int
type SortBigInts struct{ BigInts }

func (bi BigInts) Len() int           { return len(bi) }
func (bi BigInts) Swap(i, j int)      { bi[i], bi[j] = bi[j], bi[i] }
func (bi BigInts) Less(i, j int) bool { return bi[i].Cmp(bi[j]) < 0 }

func TestPowWindow(t *testing.T) {

	checkPowWindow(t, &TestElement{big.NewInt(2), *big.NewInt(100)}, big.NewInt(2), big.NewInt(4))

	checkPowWindow(t, &TestElement{big.NewInt(10), *big.NewInt(100)}, big.NewInt(10), big.NewInt(0))

	testElem := &TestElement{makeBigInt("3"), *makeBigInt("730750818665451621361119245571504901405976559617")}
	checkPowWindow(t,
		testElem,
		makeBigInt("346147755795474257120521634428450035879485727536"),
		makeBigInt("162545157220080657869228973848821629858076108602"))

	for i := 0; i < 10*1000; i++ {
		testVals := []*big.Int{makeRandBigInt(), makeRandBigInt(), makeRandBigInt()}
		sort.Sort(SortBigInts{testVals})
		// using the lowest rand as the base, next as the exponent and largest as mod
		checkPowWindowBigInt(t, &TestElement{testVals[0], *testVals[2]}, testVals[1])
	}
}

func testFrozeness(t *testing.T, x *BigInt, expect *BigInt, calc func (*BigInt) *BigInt) {
	save := x.copyUnfrozen()
	y := calc(x)
	if !save.IsEqual(x) {
		t.Errorf("Frozen value should have stayed the same: %s", save.String())
	}
	if !y.IsEqual(expect) {
		t.Errorf("Got wrong calc value: expected %s, got %s", expect.String(), y.String())
	}
}

func TestBigIntMath(t *testing.T) {

	test100 := MakeBigInt(100, false)
	test200 := MakeBigInt(200, false)

	testMod := big.NewInt(1000003) // need an odd prime or mod sqrt panics

	// expect test100 to mutate
	test100.Add(test200, testMod)
	if !test100.IsEqual(MakeBigInt(300, false)) {
		t.Errorf("Addition failed: expected 300, got %s", test100.String())
	}

	// reset to 100 - frozen
	test100 = MakeBigInt(100, true)
	testFrozeness(t, test100, MakeBigInt(200, false), func(x *BigInt) *BigInt { return x.Add(test100, testMod) } )
	testFrozeness(t, test100, MakeBigInt(0, false), func(x *BigInt) *BigInt { return x.Sub(test100, testMod) } )
	testFrozeness(t, test100, MakeBigInt(100, false), func(x *BigInt) *BigInt { return x.mod(testMod) } )
	testFrozeness(t, test100, MakeBigInt(10000, false), func(x *BigInt) *BigInt { return x.Mul(test100, testMod) } )
	testFrozeness(t, test100, MakeBigInt(10000, false), func(x *BigInt) *BigInt { return x.Square(testMod) } )
	testFrozeness(t, test100, MakeBigInt(10, false), func(x *BigInt) *BigInt { return x.sqrt(testMod) } )
	testFrozeness(t, test100, MakeBigInt(330001, false), func(x *BigInt) *BigInt { return x.invert(testMod) } )
}
