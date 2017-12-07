package field

import "math/big"

var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var TWO = big.NewInt(2)

var BI_ZERO = MakeBigInt(0)

type Field interface {
}

type Element interface {
	Copy() Element
	Mul(*Element) Element
	SetToOne() Element
}

type BaseField struct {
	LengthInBytes int
}

type PointLike interface {
	X() *BigInt
	Y() *BigInt
}

func powWindow(base *Element, exp *big.Int) *Element {

	result := (*base).Copy()
	result.SetToOne()

	if exp.Sign() == 0 {
		return &result
	}

	k := optimalPowWindowSize(exp)
	lookups := buildPowWindow(k, base)

	word := uint(0)
	wordBits := uint(0)

	inWord := false
	for s := exp.BitLen() - 1; s >= 0; s-- {
		result.Mul(&result)

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
			result.Mul(&(*lookups)[word])
			inWord = false
		}
	}

	return &result
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

func buildPowWindow(k uint, base *Element) *[]Element {

	if k < 1 {
		return nil
	}

	lookupSize := 1 << k
	lookups := make([]Element, lookupSize)

	lookups[0] = (*base).Copy().SetToOne()
	for x := 1; x < lookupSize; x++ {
		lookups[x] = lookups[x-1].Copy()
		lookups[x].Mul(base)
	}

	return &lookups
}
