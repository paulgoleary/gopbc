package field

import "math/big"

var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var TWO = big.NewInt(2)
var THREE = big.NewInt(3)

var BI_ZERO = MakeBigInt(0)
var BI_ONE = MakeBigInt(1)
var BI_THREE = MakeBigInt(3)

/*
	BigInt is intended to represent the base level of integer modular math for field computations.
	What may be a bit confusing (and I need to think about) is that I don't intend this to be a replacement for big.Int everywhere.
	The full name here is more explicit: field.BigInt - that is, a large integer that is a component of a field, which implies/requires modular math.
	Another worthy goal here might be for the wrapper to enforce some level of invariance outside the `field` package
*/
type BigInt big.Int

func MakeBigInt(x int64) *BigInt {
	return (*BigInt)(big.NewInt(x))
}

func MakeBigIntStr(x string) *BigInt {
	ret := big.Int{}
	ret.SetString(x, 10)
	return (*BigInt)(&ret)
}

func (bi *BigInt) isZero() bool {
	return (*big.Int)(bi).Cmp(ZERO) == 0
}

func (bi *BigInt) copy() *BigInt {
	if bi == nil {
		return nil
	}
	newBigInt := new(BigInt)
	(*big.Int)(newBigInt).SetBytes((*big.Int)(bi).Bytes())
	return newBigInt
}

func (bi *BigInt) setBytes(bytes []byte) {
	(*big.Int)(bi).SetBytes(bytes)
}

func (bi *BigInt) IsEqual(in *BigInt) bool {
	return (*big.Int)(bi).Cmp((*big.Int)(in)) == 0
}

func (bi *BigInt) add(in *BigInt, modIn *big.Int) *BigInt {
	(*big.Int)(bi).Add((*big.Int)(bi), (*big.Int)(in))
	return bi.mod(modIn)
}

func (bi *BigInt) sub(in *BigInt, modIn *big.Int) *BigInt {
	(*big.Int)(bi).Sub((*big.Int)(bi), (*big.Int)(in))
	return bi.mod(modIn)
}

func (bi *BigInt) mod(in *big.Int) *BigInt {
	(*big.Int)(bi).Mod((*big.Int)(bi), in)
	return bi
}

func (bi *BigInt) mul(in *BigInt, modIn *big.Int) *BigInt {
	(*big.Int)(bi).Mul((*big.Int)(bi), (*big.Int)(in))
	return bi.mod(modIn)
}

func (bi *BigInt) square(modIn *big.Int) *BigInt {
	(*big.Int)(bi).Mul((*big.Int)(bi), (*big.Int)(bi))
	return bi.mod(modIn)
}

func (bi *BigInt) invert(mod *big.Int) *BigInt {
	(*big.Int)(bi).ModInverse((*big.Int)(bi), mod)
	return bi
}

func (bi *BigInt) String() string {
	return (*big.Int)(bi).String()
}

type Field interface {
}

type Element interface {
	Copy() Element
	Mul(Element) Element
	SetToOne() Element
}

type BaseField struct {
	LengthInBytes int
}

type PointLike interface {
	X() *BigInt
	Y() *BigInt
}

func powWindow(base Element, exp *big.Int) Element {

	result := base.Copy()
	result.SetToOne()

	if exp.Sign() == 0 {
		return result
	}

	k := optimalPowWindowSize(exp)
	lookups := buildPowWindow(k, base)

	word := uint(0)
	wordBits := uint(0)

	inWord := false
	for s := exp.BitLen() - 1; s >= 0; s-- {
		result.Mul(result)

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
			result.Mul((*lookups)[word])
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

func buildPowWindow(k uint, base Element) *[]Element {

	if k < 1 {
		return nil
	}

	lookupSize := 1 << k
	lookups := make([]Element, lookupSize)

	lookups[0] = base.Copy().SetToOne()
	for x := 1; x < lookupSize; x++ {
		lookups[x] = lookups[x-1].Copy()
		lookups[x].Mul(base)
	}

	return &lookups
}
