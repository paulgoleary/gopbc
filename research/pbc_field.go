package research

import (
	"math/big"
	// "crypto/rand"
)

type Field interface {
}

/*
   public CurveField(SecureRandom random, Element a, Element b, BigInteger order, BigInteger cofac, byte[] genNoCofac) {
       super(random, (F) a.getField());

       this.random = random;
       this.a = a;
       this.b = b;
       this.order = order;
       this.cofac = cofac;

       initGen(genNoCofac);
   }
*/

type CurveField struct {
	a          ZrElement
	b          ZrElement
	order      big.Int
	cofac      big.Int
	gen        CurveElement // TODO: not sure here...
	genNoCofac CurveElement
}

type ZrField struct {
}

type CurveElement struct {
}

type ZrElement struct {
	elemField Field
	Data      []big.Int
}

func (e ZrElement) PowZn(eZn ZrElement) ZrElement {
	return e
}

var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var TWO = big.NewInt(2)

/**
This is the mod square-and-multiply algorithm with the sliding window optimization
Implemented here with base parameters - ie. independent of fields - because it only applies to scalar values
*/
func powWindow(base *big.Int, exp *big.Int, mod *big.Int) *big.Int {
	if exp.Sign() == 0 {
		return ONE
	}

	k := optimalPowWindowSize(exp)
	lookups := buildPowWindow(k, base)

	result := big.NewInt(1)

	word := uint(0)
	wordBits := uint(0)

	inWord := false
	for s := exp.BitLen() - 1; s >= 0; s-- {
		result.Mul(result, result)
		result.Mod(result, mod)

		bit := exp.Bit(s)

		if !inWord && bit == 0 {
			continue
		}

		if !inWord {
			inWord = true
			word = 1
			wordBits = 1
		} else {
			word = (word << 1) + bit
			wordBits++
		}

		if wordBits == k || s == 0 {
			result.Mul(result, &(*lookups)[word])
			result.Mod(result, mod)
			inWord = false
		}
	}

	return result
}

func optimalPowWindowSize(exp *big.Int) uint {

	expBits := exp.BitLen()

	switch {
	case expBits > 9065:
		return 8
	case expBits > 3529:
		return 7
	case expBits > 1324:
		return 6
	case expBits > 474:
		return 5
	case expBits > 157:
		return 4
	case expBits > 47:
		return 3
	default:
		return 2
	}
}

func buildPowWindow(k uint, base *big.Int) *[]big.Int {

	if k < 1 {
		return nil
	}

	lookupSize := 1 << k
	lookups := make([]big.Int, lookupSize)

	lookups[0].Set(ONE)
	for x := 1; x < lookupSize; x++ {
		lookups[x].Set(&lookups[x-1])
		lookups[x].Mul(&lookups[x], base)
	}

	return &lookups
}
