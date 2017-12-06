package field

import (
	"math/big"
	"math/rand"
	// "sort"
	"testing"
)

type TestElement struct {
	data *big.Int
	mod big.Int
}

// implement Element for TestElement

func (elem *TestElement) Copy() Element {
	newElem := new(TestElement)
	newElem.data = new(big.Int)
	newElem.data.SetBytes(elem.data.Bytes()) // need to make 'deep copy' of mutable data
	newElem.mod = elem.mod
	return newElem
}

func (elem *TestElement) Mul( mulElem *Element ) Element {
	teIn := (*mulElem).(*TestElement)
	elem.data.Mul(elem.data, teIn.data) // TODO: this is sus ...
	elem.data.Mod(elem.data, &elem.mod)
	return elem
}

func (elem *TestElement ) SetToOne() Element {
	elem.data.Set(ONE)
	return elem
}

// validate that TestElement satisfies Element interface
var _ Element = (*TestElement)(nil)

func checkPowWindowBigInt(t *testing.T, testBase *TestElement, testExp *big.Int) {
	expectedBigInt := new(big.Int).Exp(testBase.data, testExp, &testBase.mod)
	checkPowWindow(t, testBase, testExp, expectedBigInt)
}

func checkPowWindow(t *testing.T, testBase *TestElement, testExp, expectedVal *big.Int) {
	var elem Element = testBase
	testPow := powWindow(&elem, testExp)
	if (*testPow).(*TestElement).data.Cmp(expectedVal) != 0 {
		t.Errorf("powWindow exponent result was wrong, got: %d, want: %d.", testPow, expectedVal)
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

	/*
	for i := 0; i < 10*1000; i++ {
		// TODO: generator in Go?
		testVals := []*big.Int{makeRandBigInt(), makeRandBigInt(), makeRandBigInt()}
		sort.Sort(SortBigInts{testVals})
		// using the lowest rand as the base, next as the exponent and largest as mod
		checkPowWindowBigInt(t, testVals[0], testVals[1], testVals[2])
	}
	*/
}

