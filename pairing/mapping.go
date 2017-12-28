package pairing

import (
	"gobdc/field"
	"log"
	"math/big"
)

type Mapping interface  {
	pairing(var1 field.PointElement, var2 field.PointElement) field.Element
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

func (pm *TypeATateNafProjMillerPairingMap) pairing(P field.PointElement, Q field.PointElement) field.Element {

	f := pm.Fq2.MakeOne()
	// u := pm.Fq2.MakeElement()

	// JacobPoint V = new JacobPoint(P.getX(), P.getY(), P.getX().getField().newOneElement());
	V := &JacobPoint{P.X(), P.Y(), field.BI_ONE}

	nP := P.Negate()
	field.Trace(nP)

	// Element a = this.pairing.Fp.newElement();
	// Element b = this.pairing.Fp.newElement();
	// Element c = this.pairing.Fp.newElement();
	// a := pm.Fq.NewZero()
	// b := pm.Fq.NewZero()
	// c := pm.Fq.NewZero()

	// TODO !
	var a, b, c *field.BigInt
	for i := len(pm.rNAF) - 2; i >= 0; i-- {
		V, a, b, c = pm.twice(V)
		field.Trace(V, a, b, c)
		u := pm.millerStep(a, b, c, Q.X(), Q.Y())
		field.Trace(u, a, b, c)
		f = f.Square().MulPoint(u)
		field.Trace(f, a, b, c)

		switch rn := pm.rNAF[i]; rn {
		case -1, 1:
			field.Trace(V, a, b, c)
			if rn == -1 {
				V, a, b, c = pm.add(V, nP)
			} else {
				V, a, b, c = pm.add(V, P)
			}
			field.Trace(V, a, b, c)
			u = pm.millerStep(a, b, c, Q.X(), Q.Y())
			f = f.MulPoint(u)
			field.Trace(f)
		case 0: // NOP
		default:
			log.Panicf("this should not happen")
		}
	}

	/*
	for(int i = this.r.length - 2; i >= 0; --i) {
	this.twice(V, a, b, c);
	this.millerStep(u, a, b, c, Q.getX(), Q.getY());
	f.square().mul(u);
	switch(this.r[i]) {
	case -1:
	this.add(V, nP, a, b, c);
	this.millerStep(u, a, b, c, Q.getX(), Q.getY());
	f.mul(u);
	break;
	case 1:
	this.add(V, P, a, b, c);
	this.millerStep(u, a, b, c, Q.getX(), Q.getY());
	f.mul(u);
	}
	}
	*/

	field.Trace(f)
	/*
	    Point out = (Point)this.pairing.Fq2.newElement();
    this.tatePow(out, f, this.pairing.phikOnr);
    return new GTFiniteElement(this, (GTFiniteField)this.pairing.getGT(), out);
	 */

	 pm.tatePow(f, pm.phikOnr)

	return nil
}

func (pm *TypeATateNafProjMillerPairingMap) twice(V *JacobPoint) (*JacobPoint, *field.BigInt, *field.BigInt, *field.BigInt) {

	targetOrder := pm.Fq.FieldOrder

	// Element x = V.getX();
	// Element y = V.getY();
	// Element z = V.getZ();

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

	// Element t1 = y.duplicate().square();
	// Element t2 = x.duplicate().mul(t1).twice().twice();
	// b.set(z).square();
	// a.set(x).square().mul(3).add(b.duplicate().square());
	// c.set(a).mul(x).sub(t1).sub(t1);
	// z.mul(y).twice();
	// V.setX(a.duplicate().square().sub(t2.duplicate().twice()));
	// V.setY( a.duplicate().mul(t2.duplicate().sub(V.getX())) .sub(t1.duplicate().square().twice().twice().twice()) );
	// a.mul(b);
	// b.mul(z);

	V.freeze()
	a.Freeze()
	b.Freeze()
	c.Freeze()
	return V, a, b, c
}

// we mutate `out` in place here
func (pm *TypeATateNafProjMillerPairingMap) millerStep(a, b, c, Qx, Qy *field.BigInt) field.PointElement {
	a.Freeze()
	b.Freeze()
	c.Freeze()
	targetOrder := pm.Fq.FieldOrder
	x := c.Add(a.Mul(Qx, targetOrder), targetOrder)
	y := b.Mul(Qy, targetOrder)
	return pm.Fq2.MakeElement(x, y)

	// out.getX().set(c).add(a.duplicate().mul(Qx));
	// out.getY().set(b).mul(Qy);
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

/*
  final void tatePow(Point out, Point in, BigInteger cofactor) {
    Element in1 = in.getY();
    Point temp = (Point)in.duplicate().invert();
    in1.negate();
    in.mul(temp);
    this.lucasOdd(out, in, temp, cofactor);
  }
 */

func (pm *TypeATateNafProjMillerPairingMap) tatePow(in field.PointElement, cofactor *big.Int) field.PointElement {

	// targetOrder := pm.Fq.FieldOrder

	// Element in1 = in.getY();
	// in1.negate();
	// inY := in.Y().Negate(targetOrder)
	// TODO: not used? or expected to get picked up as a side-effect?

	// Point temp = (Point)in.duplicate().invert();
	tempPoint := in.Invert()

	in = in.MulPoint(tempPoint)

	// this.lucasOdd(out, in, temp, cofactor);

	return nil
}
