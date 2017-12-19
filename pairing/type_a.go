package pairing

import (
	"fmt"
	"gobdc/field"
	"math/big"
	"log"
)

type TypeAPairingParams struct {
	exp2            int
	exp1            int
	sign1           int
	sign0           int
	r               *big.Int // r = 2^exp2 + sign1 * 2^exp1 + sign0 * 1
	q               *big.Int // we work in E(F_q) (and E(F_q^2))
	h               *big.Int // r * h = q + 1
	genNoCofacBytes *[]byte
}

type TypeAPairing struct {
	TypeAPairingParams
	BasePairing
	Fq *field.ZrField
	Fq2 *field.D2ExtensionQuadField
}

func (pairing *TypeAPairing) initTypeAPairingParams(params *PairingParameters) {
	pairingType := (*params)["type"]
	if pairingType != aType {
		panic(fmt.Sprintf("Invalid pairing type '%s' - expected 'a'", pairingType))
	}

	pairing.exp2 = params.getInt("exp2")
	pairing.exp1 = params.getInt("exp1")
	pairing.sign1 = params.getInt("sign1")

	pairing.r = params.getBigInt("r") // r = 2^exp2 + sign1 * 2^exp1 + sign0 * 1
	pairing.q = params.getBigInt("q") // we work in E(F_q) (and E(F_q^2))
	pairing.h = params.getBigInt("h") // r * h = q + 1

	// TODO: need to test that this decodes in the same way as PBC, jPBC, etc.
	pairing.genNoCofacBytes = params.getBytes("genNoCofac", nil)
}

// TODO: compatibility with jPBC and PBC ???
const (
	NAF_MILLER_PROJECTIVE_METHOD = "naf-miller-projective"
)

func (pairing *TypeAPairing) initTypeAPairingMap(params *PairingParameters) {
	method := params.getString("method", NAF_MILLER_PROJECTIVE_METHOD)
	if method != NAF_MILLER_PROJECTIVE_METHOD {
		log.Panicf("Pairing method currently unsupported: %s", method)
	}

	pairing.mapping = MakeTypeATateNafProjMillerPairingMap(pairing)



}

/*
   protected Field initFp(BigInteger order) {
       return new ZrField(random, order);
   }

   protected Field<? extends Point> initEq() {
       // Remember the curve is: y^2 = x^3 + ax
       return new CurveField<Field>(random,
                                    Fq.newOneElement(),   // a
                                    Fq.newZeroElement(),  // b
                                    r,                    // order
                                    h,                    // cofactor  (r*h)=q+1=#E(F_q)
                                    genNoCofac);
   }
*/

func (pairing *TypeAPairing) makeEq() *field.CurveField {
	return field.MakeCurveField(
		pairing.Fq.NewOneElement(),
		pairing.Fq.NewZeroElement(),
		pairing.r,
		pairing.h,
		pairing.genNoCofacBytes)
}

func (pairing *TypeAPairing) initTypeAPairingFields(params *PairingParameters) {
	// Init Zr
	pairing.Zr = field.MakeZrField(pairing.r)

	// Init Fq
	pairing.Fq = field.MakeZrField(pairing.q)

	// TODO: any reason to have Eq ?
	// Init Eq
	// pairing.Eq = initEq();

	pairing.Fq2 = field.MakeD2ExtensionQuadField(pairing.Fq)

	// k=2, hence phi_k(q) = q + 1, phikOnr = (q+1)/r
	// phikOnr = h;

	// Init G1, G2, GT
	pairing.G1 = pairing.makeEq()
	pairing.G2 = pairing.G1
	// GT = initGT();
}

func MakeTypeAPairing(params *PairingParameters) *TypeAPairing {
	pairing := new(TypeAPairing)
	pairing.initTypeAPairingParams(params)
	pairing.initTypeAPairingFields(params)
	pairing.initTypeAPairingMap(params)
	return pairing
}
