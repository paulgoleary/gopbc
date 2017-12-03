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
	a          ZrElement
	b          ZrElement
	order      big.Int
	cofac      big.Int
	gen        CurveElement // TODO: not sure here...
	genNoCofac CurveElement
}

type CurveElement struct {
	ElemField CurveField
	DataX     big.Int
	DataY     big.Int
}

// CurveField

func MakeCurveField() CurveField {
	field := new(CurveField)

	return *field
}