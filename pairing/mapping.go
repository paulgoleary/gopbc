package pairing

import "gobdc/field"

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

	// Point f = (Point)this.pairing.Fq2.newOneElement();
	// Point u = (Point)this.pairing.Fq2.newElement();
	f := pm.Fq2.MakeElement().SetToOne()
	u := pm.Fq2.MakeElement()

	// JacobPoint V = new JacobPoint(P.getX(), P.getY(), P.getX().getField().newOneElement());
	V := JacobPoint{P.X(), P.Y(), field.BI_ONE}

	// Point nP = (Point)P.duplicate().negate()
	nP := P.Negate()

	// Element a = this.pairing.Fp.newElement();
	// Element b = this.pairing.Fp.newElement();
	// Element c = this.pairing.Fp.newElement();
	a := pm.Fq.NewZeroElement()
	b := pm.Fq.NewZeroElement()
	c := pm.Fq.NewZeroElement()

	// TODO !
	for i := len(pm.rNAF) - 2; i >= 0; i-- {

	}

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


	return nil
}
