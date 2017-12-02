package research

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
)

type PairingParameters map[string]string

func (params PairingParameters) getInt(paramName string) int {
	ret, err := strconv.Atoi(params[paramName])
	if err != nil {
		panic(err.Error()) // TODO: additional error context?
	}
	return ret
}

// TODO: error handling
func (params PairingParameters) getBigInt(paramName string) *big.Int {
	paramVal := params[paramName]
	ret := new(big.Int)
	ret.SetString(paramVal, 10)
	return ret
}

func (params PairingParameters) getBytes(paramName string, defaultVal *[]byte) *[]byte {
	paramVal := params[paramName]
	if paramVal == "" {
		return defaultVal
	}
	decoded, err := base64.StdEncoding.DecodeString(paramVal)
	if err != nil {
		panic(fmt.Sprintf("Could not base64 decode param name '%s', value: %s", paramName, paramVal))
	}
	return &decoded
}

func (params PairingParameters) getString(paramName string, defaultVal string) string {
	paramVal := params[paramName]
	if paramVal == "" {
		return defaultVal
	}
	return paramVal
}

func getCompatParams() (params PairingParameters) {

	/*
	   *** this is now actually *also* cribbed from the bdc project ...
	   taken from PBC-provided param file 'a.param'
	   type a
	   q 8780710799663312522437781984754049815806883199414208211028653399266475630880222957078625179422662221423155858769582317459277713367317481324925129998224791
	   h 12016012264891146079388821366740534204802954401251311822919615131047207289359704531102844802183906537786776
	   r 730750818665451621361119245571504901405976559617
	   exp2 159
	   exp1 107
	   sign1 1
	   sign0 1
	*/

	params = PairingParameters{}
	params["type"] = "a"
	params["q"] = "8780710799663312522437781984754049815806883199414208211028653399266475630880222957078625179422662221423155858769582317459277713367317481324925129998224791"
	params["r"] = "730750818665451621361119245571504901405976559617"
	params["h"] = "12016012264891146079388821366740534204802954401251311822919615131047207289359704531102844802183906537786776"
	params["exp1"] = "107"
	params["exp2"] = "159"
	params["sign0"] = "1"
	params["sign1"] = "1"

	params["genNoCofac"] = "WLeuxaO0DxaW+oJ4vrLKgkq91prZNLGQUVoXH4gIx6AGIS7vrU7Fq3/5DfYTRHfpnOCIuo96hfRwTzUTf2+EUndlGtVaI05vjWxsIaCqKSPtq+xYpr0jaGVVwnXojhjbi0AeR/JvjiIaF9WFjSRzqEvR1WHp0LkJRrtBfNcA0k4="

	return
}

const (
	aType = "a"
)

type TypeAPairingParams struct {
	exp2       int
	exp1       int
	sign1      int
	sign0      int
	r          *big.Int // r = 2^exp2 + sign1 * 2^exp1 + sign0 * 1
	q          *big.Int // we work in E(F_q) (and E(F_q^2))
	h          *big.Int // r * h = q + 1
	genNoCofac *[]byte
}

// TODO!
type BasePairing struct {
	// protected SecureRandom random;
	G1, G2, GT, Zr Field
	// protected PairingMap pairingMap;
}

type TypeAPairing struct {
	TypeAPairingParams
	BasePairing
}

func (pairing *TypeAPairing) initTypeAPairingParams(params PairingParameters) {
	pairingType := params["type"]
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
	pairing.genNoCofac = params.getBytes("genNoCofac", nil)
}

// TODO: compatibility with jPBC and PBC ???
const (
	NAF_MILLER_PROJECTTIVE_METHOD = "naf-miller-projective"
)

func (pairing *TypeAPairing) initTypeAPairingMap(params PairingParameters) {
	method := params.getString("method", NAF_MILLER_PROJECTTIVE_METHOD)
	println(fmt.Sprintf("CURRENTLY NOT IMPLEMENTED!: %s", method)) // TODO!!!
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

// TODO!!!
func (pairing *TypeAPairing) makeEq() *CurveField {
	return new(CurveField)
}

func (pairing *TypeAPairing) initTypeAPairingFields(params PairingParameters) {
	// Init Zr
	// Zr = initFp(r);

	// Init Fq
	// Fq = initFp(q);

	// TODO: any reason to have Eq ?
	// Init Eq
	// pairing.Eq = initEq();

	// Init Fq2
	// Fq2 = initFi();

	// k=2, hence phi_k(q) = q + 1, phikOnr = (q+1)/r
	// phikOnr = h;

	// Init G1, G2, GT
	pairing.G1 = pairing.makeEq()
	pairing.G2 = pairing.G1
	// GT = initGT();
}

func MakeTypeAPairing(params PairingParameters) *TypeAPairing {
	pairing := new(TypeAPairing)
	pairing.initTypeAPairingParams(params)
	pairing.initTypeAPairingMap(params)
	pairing.initTypeAPairingFields(params)
	return pairing
}
