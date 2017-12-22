package pairing

import (
	"gobdc/field"
	"log"
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

// TODO !!!
func (pm *TypeATateNafProjMillerPairingMap) pairing(P field.PointElement, Q field.PointElement) field.Element {

	// TODO: maybe an explicit make for one element so we don't need to pass in nil's ?
	f := pm.Fq2.MakeElement(nil, nil).SetToOne()
	// u := pm.Fq2.MakeElement()

	// JacobPoint V = new JacobPoint(P.getX(), P.getY(), P.getX().getField().newOneElement());
	V := &JacobPoint{field.PointLike{P.X(), P.Y()}, field.BI_ONE}

	nP := P.Negate()

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
		f = f.Square().Mul(u)
		field.Trace(f, a, b, c)

		switch rn := pm.rNAF[i]; rn {
		case -1, 1:
			if rn == -1 {
				pm.add(V, nP, a, b, c)
			} else {
				pm.add(V, P, a, b, c)
			}
			u = pm.millerStep(a, b, c, Q.X(), Q.Y())
			f = f.Mul(u)
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

	return nil
}

func (pm *TypeATateNafProjMillerPairingMap) twice(V *JacobPoint) (*JacobPoint, *field.BigInt, *field.BigInt, *field.BigInt) {

	targetOrder := pm.Fq.FieldOrder

	// Element x = V.getX();
	// Element y = V.getY();
	// Element z = V.getZ();

	// TODO: validate frozen? why not just just freeze? can't hurt ...
	t1 := V.DataY.Square(targetOrder)
	t1.Freeze()
	t2 := V.DataX.Mul(t1, targetOrder).Mul(field.BI_FOUR, targetOrder)
	t2.Freeze()

	b := V.z.Square(targetOrder)
	b.Freeze()
	a := V.DataX.Square(targetOrder).Mul(field.BI_THREE, targetOrder).Add(b.Square(targetOrder), targetOrder)
	a.Freeze()
	c := a.Mul(V.DataX, targetOrder).Sub(t1, targetOrder).Sub(t1, targetOrder)
	V.z = V.z.Mul(V.DataY, targetOrder).Mul(field.BI_TWO, targetOrder)
	V.DataX = a.Square(targetOrder).Sub(t2.Mul(field.BI_TWO, targetOrder), targetOrder)
	V.DataY = a.Mul(t2.Sub(V.DataX, targetOrder), targetOrder).Sub(t1.Square(targetOrder).Mul(field.BI_EIGHT, targetOrder), targetOrder)
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

// final void add(JacobPoint V, Point P, Element a, Element b, Element c) {
func (pm *TypeATateNafProjMillerPairingMap) add(V *JacobPoint, P field.PointElement , a *field.BigInt, b *field.BigInt, c *field.BigInt) {
	// TODO !!!
}
