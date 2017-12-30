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

	V := &JacobPoint{P.X(), P.Y(), field.BI_ONE}

	nP := P.NegateY()

	var a, b, c *field.BigInt
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

func (pm *TypeATateNafProjMillerPairingMap) twice(V *JacobPoint) (*JacobPoint, *field.BigInt, *field.BigInt, *field.BigInt) {

	targetOrder := pm.Fq.FieldOrder

	// TODO: validate frozen? why not just just freeze? can't hurt ...
	t1 := V.y.Square(targetOrder)
	t1.Freeze()
	t2 := V.x.Mul(t1, targetOrder).Mul(field.BI_FOUR, targetOrder)
	t2.Freeze()

	b := V.z.Square(targetOrder)
	b.Freeze()
	a := V.x.Square(targetOrder).Mul(field.BI_THREE, targetOrder).Add(b.Square(targetOrder), targetOrder)
	a.Freeze()
	c := a.Mul(V.x, targetOrder).Sub(t1, targetOrder).Sub(t1, targetOrder)
	V.z = V.z.Mul(V.y, targetOrder).Mul(field.BI_TWO, targetOrder)
	V.x = a.Square(targetOrder).Sub(t2.Mul(field.BI_TWO, targetOrder), targetOrder)
	V.y = a.Mul(t2.Sub(V.x, targetOrder), targetOrder).Sub(t1.Square(targetOrder).Mul(field.BI_EIGHT, targetOrder), targetOrder)
	a = a.Mul(b, targetOrder)
	b = b.Mul(V.z, targetOrder)

	V.freeze()
	a.Freeze()
	b.Freeze()
	c.Freeze()
	return V, a, b, c
}

func (pm *TypeATateNafProjMillerPairingMap) millerStep(a, b, c, Qx, Qy *field.BigInt) field.PointElement {
	a.Freeze()
	b.Freeze()
	c.Freeze()
	targetOrder := pm.Fq.FieldOrder
	x := c.Add(a.Mul(Qx, targetOrder), targetOrder)
	y := b.Mul(Qy, targetOrder)
	return pm.Fq2.MakeElement(x, y)
}

func (pm *TypeATateNafProjMillerPairingMap) add(V *JacobPoint, P field.PointElement) (*JacobPoint, *field.BigInt, *field.BigInt, *field.BigInt) {

	targetOrder := pm.Fq.FieldOrder
	// V is used, mutated and returned
	// P is used but *not* returned
	// a, b, c are returned but *not* used

	V.freeze() // make sure V is frozen
	// P should be frozen because it's an Element

    t1 := V.z.Square(targetOrder)
    t2 := V.z.Mul(t1, targetOrder)
    t3 := P.X().Mul(t1, targetOrder)
    t3.Freeze()
    t4 := P.Y().Mul(t2, targetOrder)
    t4.Freeze()
    t5 := t3.Sub(V.x, targetOrder)
    t5.Freeze()
    t6 := t4.Sub(V.y, targetOrder)
    t6.Freeze()
    t7 := t5.Square(targetOrder)
    t8 := t5.Mul(t7, targetOrder)
    t8.Freeze()
    t9 := V.x.Mul(t7, targetOrder)
    t9.Freeze()
    x3 := t6.Square(targetOrder).Sub(t8.Add(t9.Mul(field.BI_TWO, targetOrder), targetOrder), targetOrder)
    y3 := t6.Mul(t9.Sub(x3, targetOrder), targetOrder).Sub(V.y.Mul(t8, targetOrder), targetOrder)
    z3 := V.z.Mul(t5, targetOrder)

    V.x = x3
    V.y = y3
    V.z = z3
    V.freeze()

    a := t6
    b := z3
    c := t6.Mul( P.X(), targetOrder).Sub(z3.Mul(P.Y(), targetOrder), targetOrder) // z3 is frozen through V

    return V, a, b, c
}

func (pm *TypeATateNafProjMillerPairingMap) tatePow(in field.PointElement, cofactor *big.Int) field.PointElement {

	tempPoint := in.Invert()
	in = in.NegateY().MulPoint(tempPoint)
	return pm.lucasOdd(in, tempPoint, cofactor)
}

func (pm *TypeATateNafProjMillerPairingMap) lucasOdd(in field.PointElement, temp field.PointElement, cofactor *big.Int) field.PointElement {

	targetOrder := pm.Fq.FieldOrder

	t0 := field.BI_TWO
	t1 := in.X().Mul(field.BI_TWO, targetOrder)

	// TODO: frozen-ness?
	// no need to Copy t0 since it is frozen - it will just be copied again
	v0 := t0
	v1 := t1.Copy()

	for j := cofactor.BitLen() - 1; j != 0; j-- {
		if cofactor.Bit(j) != 0 {
			v0 = v0.Mul(v1, targetOrder).Sub(t1, targetOrder)
			v1 = v1.Square(targetOrder).Sub(t0, targetOrder)
		} else {
			v1 = v1.Mul(v0, targetOrder).Sub(t1, targetOrder)
			v0 = v0.Square(targetOrder).Sub(t0, targetOrder)
		}
	}

	v1 = v1.Mul(v0, targetOrder).Sub(t1, targetOrder)
	v0 = v0.Square(targetOrder).Sub(t0, targetOrder)

	v0.Freeze() // sneaky mutation below...
	v1 = v1.Mul(field.BI_TWO, targetOrder).Sub(v0.Mul(t1, targetOrder), targetOrder)
	t1 = t1.Square(targetOrder).Sub(t0, targetOrder).Sub(t0, targetOrder)
	v1 = v1.Mul(t1.Invert(targetOrder), targetOrder)

	v0 = v0.Mul(field.BI_TWO.Invert(targetOrder), targetOrder) // TODO jPBC pre-computes the mod inverse of two and uses it for this calc
	v1 = v1.Mul(in.Y(), targetOrder)

	return pm.Fq2.MakeElement(v0, v1)
}
