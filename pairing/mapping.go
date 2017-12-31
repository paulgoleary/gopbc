package pairing

import (
	"gobdc/field"
	"log"
	"math/big"
)

type Mapping interface  {
	pairing(var1 field.PointElement, var2 field.PointElement) *GTFiniteElement
	// isProductPairingSupported() bool
	// pairings(var1 []field.Element, var2 []field.Element ) field.Element
	// finalPow(var1 field.Element)
	// isAlmostCoddh(var1 field.Element, var2 field.Element, var3 field.Element, var4 field.Element) bool
	// getPairingPreProcessingLengthInBytes() int
	// pairingPoint(var1 field.PointLike) PreProcessing
	// pairingBytes(var1 []byte, var2 int) PreProcessing
}

type TypeATateNafProjMillerPairingMap struct {
	TypeAPairing
	rNAF []int8
}

func MakeTypeATateNafProjMillerPairingMap(pairing *TypeAPairing) *TypeATateNafProjMillerPairingMap {

	pairingMap := new(TypeATateNafProjMillerPairingMap)
	pairingMap.TypeAPairing = *pairing
	pairingMap.rNAF = field.NAF(pairing.r, 2)
	return pairingMap
}

// MakeTateNafProjMillerPairingMap

var _ Mapping = (*TypeATateNafProjMillerPairingMap)(nil)

func (pm *TypeATateNafProjMillerPairingMap) pairing(P field.PointElement, Q field.PointElement) *GTFiniteElement {

	f := pm.Fq2.MakeOne()

	V := &JacobPoint{P.X(), P.Y(), pm.Fq.NewOneElement().ModInt}

	nP := P.NegateY()

	var a, b, c *field.ModInt
	for i := len(pm.rNAF) - 2; i >= 0; i-- {
		V, a, b, c = pm.twice(V)
		u := pm.millerStep(a, b, c, Q.X(), Q.Y())
		f = f.Square().MulPoint(u)

		switch rn := pm.rNAF[i]; rn {
		case -1, 1:
			if rn == -1 {
				V, a, b, c = pm.add(V, nP)
			} else {
				V, a, b, c = pm.add(V, P)
			}
			u = pm.millerStep(a, b, c, Q.X(), Q.Y())
			f = f.MulPoint(u)
		case 0: // NOP
		default:
			log.Panicf("this should not happen")
		}
	}

	out := pm.tatePow(f, pm.phikOnr)

	return pm.GT.MakeElement(out, pm)
}

func (pm *TypeATateNafProjMillerPairingMap) twice(V *JacobPoint) (*JacobPoint, *field.ModInt, *field.ModInt, *field.ModInt) {

	// TODO: validate frozen? why not just just freeze? can't hurt ...
	t1 := V.y.Square()
	t1.Freeze()
	t2 := V.x.Mul(t1).Mul(field.MI_FOUR)
	t2.Freeze()

	b := V.z.Square()
	b.Freeze()
	a := V.x.Square().Mul(field.MI_THREE).Add(b.Square())
	a.Freeze()
	c := a.Mul(V.x).Sub(t1).Sub(t1)
	V.z = V.z.Mul(V.y).Mul(field.MI_TWO)
	V.x = a.Square().Sub(t2.Mul(field.MI_TWO))
	V.y = a.Mul(t2.Sub(V.x)).Sub(t1.Square().Mul(field.MI_EIGHT))
	a = a.Mul(b)
	b = b.Mul(V.z)

	V.freeze()
	a.Freeze()
	b.Freeze()
	c.Freeze()
	return V, a, b, c
}

func (pm *TypeATateNafProjMillerPairingMap) millerStep(a, b, c, Qx, Qy *field.ModInt) field.PointElement {
	a.Freeze()
	b.Freeze()
	c.Freeze()
	x := c.Add(a.Mul(Qx))
	y := b.Mul(Qy)
	return pm.Fq2.MakeElement(x, y)
}

func (pm *TypeATateNafProjMillerPairingMap) add(V *JacobPoint, P field.PointElement) (*JacobPoint, *field.ModInt, *field.ModInt, *field.ModInt) {

	// V is used, mutated and returned
	// P is used but *not* returned
	// a, b, c are returned but *not* used

	V.freeze() // make sure V is frozen
	// P should be frozen because it's an Element

    t1 := V.z.Square()
    t2 := V.z.Mul(t1)
    t3 := P.X().Mul(t1)
    t3.Freeze()
    t4 := P.Y().Mul(t2)
    t4.Freeze()
    t5 := t3.Sub(V.x)
    t5.Freeze()
    t6 := t4.Sub(V.y)
    t6.Freeze()
    t7 := t5.Square()
    t8 := t5.Mul(t7)
    t8.Freeze()
    t9 := V.x.Mul(t7)
    t9.Freeze()
    x3 := t6.Square().Sub(t8.Add(t9.Mul(field.MI_TWO)))
    y3 := t6.Mul(t9.Sub(x3)).Sub(V.y.Mul(t8))
    z3 := V.z.Mul(t5)

    V.x = x3
    V.y = y3
    V.z = z3
    V.freeze()

    a := t6
    b := z3
    c := t6.Mul( P.X()).Sub(z3.Mul(P.Y())) // z3 is frozen through V

    return V, a, b, c
}

func (pm *TypeATateNafProjMillerPairingMap) tatePow(in field.PointElement, cofactor *big.Int) field.PointElement {

	tempPoint := in.Invert()
	in = in.NegateY().MulPoint(tempPoint)
	return pm.lucasOdd(in, tempPoint, cofactor)
}

func (pm *TypeATateNafProjMillerPairingMap) lucasOdd(in field.PointElement, temp field.PointElement, cofactor *big.Int) field.PointElement {

	t0 := pm.Fq.NewElement(field.TWO).ModInt
	t1 := in.X().Mul(field.MI_TWO)

	// TODO: frozen-ness?
	// no need to Copy t0 since it is frozen - it will just be copied again
	v0 := t0
	v1 := t1.Copy()

	for j := cofactor.BitLen() - 1; j != 0; j-- {
		if cofactor.Bit(j) != 0 {
			v0 = v0.Mul(v1).Sub(t1)
			v1 = v1.Square().Sub(t0)
		} else {
			v1 = v1.Mul(v0).Sub(t1)
			v0 = v0.Square().Sub(t0)
		}
	}

	v1 = v1.Mul(v0).Sub(t1)
	v0 = v0.Square().Sub(t0)

	v0.Freeze() // sneaky mutation below...
	v1 = v1.Mul(field.MI_TWO).Sub(v0.Mul(t1))
	t1 = t1.Square().Sub(t0).Sub(t0)
	v1 = v1.Mul(t1.Invert())

	v0 = v0.Mul(pm.Fq.TwoInverse)
	v1 = v1.Mul(in.Y())

	return pm.Fq2.MakeElement(v0, v1)
}
