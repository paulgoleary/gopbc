package research

import (
	"testing"
	"math/big"
	"math/rand"
	"sort"
)

func checkPowWindowBigInt(t *testing.T, testBase, testExp, testMod *big.Int) {
	expectedBigInt := new(big.Int).Exp(testBase, testExp, testMod)
	checkPowWindow(t, testBase, testExp, testMod, expectedBigInt)
}

func checkPowWindow(t *testing.T, testBase, testExp, testMod, expectedVal *big.Int) {
	testPow := powWindow(testBase, testExp, testMod)
	if testPow.Cmp(expectedVal) != 0 {
		t.Errorf("powWindow exponent result was wrong, got: %d, want: %d.", testPow, expectedVal)
	}
}

func makeBigInt( valStr string ) (res *big.Int) {
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
type SortBigInts struct { BigInts }

func (bi BigInts) Len() int      { return len(bi) }
func (bi BigInts) Swap(i, j int) { bi[i], bi[j] = bi[j], bi[i] }
func (bi BigInts ) Less(i, j int) bool { return bi[i].Cmp(bi[j]) < 0 }

func TestPowWindow(t *testing.T) {

	checkPowWindow(t, big.NewInt(2), big.NewInt(2), big.NewInt(100), big.NewInt(4))

	checkPowWindow(t, big.NewInt(10), big.NewInt(10), big.NewInt(100), big.NewInt(0))

	checkPowWindow(t,
		makeBigInt("3"),
		makeBigInt("346147755795474257120521634428450035879485727536"),
		makeBigInt("730750818665451621361119245571504901405976559617"),
		makeBigInt("162545157220080657869228973848821629858076108602") )

	for i := 0; i < 10 * 1000; i++ {
		// TODO: generator in Go?
		testVals := []*big.Int{makeRandBigInt(), makeRandBigInt(), makeRandBigInt()}
		sort.Sort(SortBigInts{testVals})
		// using the lowest rand as the base, next as the exponent and largest as mod
		checkPowWindowBigInt(t, testVals[0], testVals[1], testVals[2])
	}
}