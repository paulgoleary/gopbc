package field

import (
	"math/big"
	"fmt"
)

var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var TWO = big.NewInt(2)
var THREE = big.NewInt(3)

var BI_ZERO = MakeBigInt(0, true)
var BI_ONE = MakeBigInt(1, true)
var BI_TWO = MakeBigInt(2, true)
var BI_THREE = MakeBigInt(3, true)
var BI_FOUR = MakeBigInt(4, true)
var BI_EIGHT = MakeBigInt(8, true)

/*
	BigInt is intended to represent the base level of integer modular math for field computations.
	What may be a bit confusing (and I need to think about) is that I don't intend this to be a replacement for big.Int everywhere.
	The full name here is more explicit: field.BigInt - that is, a large integer that is a component of a field, which implies/requires modular math.
*/
type BigInt struct {
	v big.Int
	frozen bool
}

func MakeBigInt(x int64, frozen bool) *BigInt {
	return &BigInt{ *big.NewInt(x), frozen  }
}

func MakeBigIntStr(x string, frozen bool) *BigInt {
	ret := big.Int{}
	ret.SetString(x, 10)
	return &BigInt{ret, frozen}
}

func (bi *BigInt) Freeze() {
	bi.frozen = true
	return
}

func (bi *BigInt) isZero() bool {
	return bi.v.Cmp(ZERO) == 0
}

func CopyFrom(bi *big.Int, frozen bool) *BigInt {
	if bi == nil {
		return nil
	}
	newBigInt := new(BigInt)
	newBigInt.v.SetBytes(bi.Bytes())
	newBigInt.frozen = frozen
	return newBigInt
}

func (bi *BigInt) copyUnfrozen() *BigInt {
	if bi == nil {
		return nil
	}
	return CopyFrom(&bi.v, false)
}

func (bi *BigInt) copy() *BigInt {
	if bi == nil {
		return nil
	}
	return CopyFrom(&bi.v, bi.frozen)
}

func (bi *BigInt) setBytes(bytes []byte) {
	bi.v.SetBytes(bytes)
}

// TODO: how do we want these functions to behave WRT nil?
func (bi *BigInt) IsEqual(in *BigInt) bool {
	if bi == nil || in == nil {
		return false
	}
	return bi.v.Cmp(&in.v) == 0
}

func (bi *BigInt) Add(in *BigInt, modIn *big.Int) *BigInt {
	if bi.frozen {
		bi = bi.copyUnfrozen()
	}
	bi.v.Add(&bi.v, &in.v)
	return bi.mod(modIn)
}

func (bi *BigInt) Sub(in *BigInt, modIn *big.Int) *BigInt {
	if bi.frozen {
		bi = bi.copyUnfrozen()
	}
	bi.v.Sub(&bi.v, &in.v)
	return bi.mod(modIn)
}

func (bi *BigInt) mod(in *big.Int) *BigInt {
	if bi.frozen {
		bi = bi.copyUnfrozen()
	}
	bi.v.Mod(&bi.v, in)
	return bi
}

func (bi *BigInt) Mul(in *BigInt, modIn *big.Int) *BigInt {
	if bi.frozen {
		bi = bi.copyUnfrozen()
	}
	bi.v.Mul(&bi.v, &in.v)
	return bi.mod(modIn)
}

func (bi *BigInt) Square(modIn *big.Int) *BigInt {
	if bi.frozen {
		bi = bi.copyUnfrozen()
	}
	bi.v.Mul(&bi.v, &bi.v)
	return bi.mod(modIn)
}

func (bi *BigInt) sqrt(modIn *big.Int) *BigInt {
	// Int.ModSqrt implements  Tonelli-Shanks and also a more optimal version when modIn = 3 mod 4
	// UGH! need to work around this bug: https://github.com/golang/go/issues/22265
	// for now always copy
	calc := bi.copyUnfrozen()
	calc.v.ModSqrt(&bi.v, modIn)
	return calc
}

func (bi *BigInt) invert(mod *big.Int) *BigInt {
	if bi.frozen {
		bi = bi.copyUnfrozen()
	}
	bi.v.ModInverse(&bi.v, mod)
	return bi
}

func (bi *BigInt) negate(modIn *big.Int) *BigInt {
	if bi.isZero() {
		return BI_ONE.copyUnfrozen()
	}
	return CopyFrom(modIn, false).Sub(bi, modIn)
}

func (bi *BigInt) String() string {
	return bi.v.String()
}

type Field interface {}

type PointField interface {
	MakeElement() PointElement
}

type Element interface {
	String() string
	Copy() Element
	Mul(Element) Element
	SetToOne() Element
}

type PointElement interface {
	Element
	X() *BigInt
	Y() *BigInt
	Negate() PointElement // TODO: might need to promote to Element?
}

type MakeElementFunc func() Element
type BaseField struct {
	LengthInBytes int
	FieldOrder *big.Int
}

type PointLike struct {
	DataX *BigInt
	DataY *BigInt
}

func (p *PointLike) String() string {
	return fmt.Sprintf("[%s],[%s]", p.DataX.String(), p.DataY.String())
}

func (p *PointLike) freeze() {
	p.DataX.Freeze()
	p.DataY.Freeze()
}

func (p *PointLike) frozen() bool {
	return p.DataX.frozen && p.DataY.frozen
}

func (p *PointLike) X() *BigInt {
	return p.DataX
}

func (p *PointLike) Y() *BigInt {
	return p.DataY
}

func powWindow(base Element, exp *big.Int) Element {

	// note: does not mutate base
	result := base.SetToOne()

	if exp.Sign() == 0 {
		return result
	}

	k := optimalPowWindowSize(exp)
	lookups := buildPowWindow(k, base)

	word := uint(0)
	wordBits := uint(0)

	inWord := false
	for s := exp.BitLen() - 1; s >= 0; s-- {
		result = result.Mul(result)

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
			result = result.Mul((*lookups)[word])
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

	// SetToOne copies ...
	lookups[0] = base.SetToOne()
	for x := 1; x < lookupSize; x++ {
		newLookup := lookups[x-1].Copy()
		lookups[x] = newLookup.Mul(base)
	}

	return &lookups
}
