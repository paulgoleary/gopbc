package field

import "math/big"

/*
   public CurveField(SecureRandom random, Element a, Element b, BigInteger order, BigInteger cofac, byte[] genNoCofac) {
       super(random, (F) a.getField());

       this.random = random;
       this.a = a;
       this.b = b;
       this.order = order;
       this.cofac = cofac;

       initGen(genNoCofac);
   }
*/

type CurveField struct {
	a          *ZrElement
	b          *ZrElement
	order      *big.Int
	cofactor   *big.Int
	gen        *CurveElement // TODO: not sure here...
	genNoCofac *CurveElement
}

type CurveElement struct {
	ElemField *CurveField
	DataX     big.Int
	DataY     big.Int
}

// CurveField

// TODO: JPBC (PBC?) handles case w/o bytes and cofactor
func (field *CurveField) initGenFromBytes( genNoCofac *[]byte ) {

	// TODO !!!
	// field.genNoCofac = field.newElementFromBytes(genNoCofac);
	// field.gen = field.genNoCofac.duplicate().mul(this.cofac);
}

/*
        // Remember the curve is: y^2 = x^3 + ax
        return new CurveField<Field>(random,
                                     Fq.newOneElement(),   // a
                                     Fq.newZeroElement(),  // b
                                     r,                    // order
                                     h,                    // cofactor  (r*h)=q+1=#E(F_q)
                                     genNoCofac);
 */
func MakeCurveField(
	a *ZrElement,
	b *ZrElement,
	order *big.Int,
	cofactor *big.Int,
	genNoCofacBytes *[]byte ) *CurveField {

	field := new(CurveField)
	field.a = a
	field.b = b
	field.order = order
	field.cofactor = cofactor
	field.initGenFromBytes(genNoCofacBytes)

	return field
}