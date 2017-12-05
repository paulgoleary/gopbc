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
	DataX     *big.Int // TODO: perhaps X and Y should be elements of the target field, as in PBC/JPBC?
	DataY     *big.Int
}

// CurveField

// TODO: JPBC (PBC?) handles case w/o bytes and cofactor
func (field *CurveField) initGenFromBytes( genNoCofac *[]byte ) {
	field.genNoCofac = field.newElementFromBytes(genNoCofac)
	field.gen = field.genNoCofac.Copy().Mul(field.cofactor)
}

func (field *CurveField) GetGen() *CurveElement {
	return field.gen
}

func (field *CurveField) getTargetField() *ZrField {
	return field.a.ElemField
}

func (field *CurveField) newElementFromBytes( elemBytes *[]byte ) *CurveElement {

	elem := CurveElement{ElemField: field}

	xBytes := (*elemBytes)[:field.getTargetField().LengthInBytes]
	yBytes := (*elemBytes)[field.getTargetField().LengthInBytes:]

	elem.DataX = new(big.Int)
	elem.DataX.SetBytes(xBytes)

	elem.DataY = new(big.Int)
	elem.DataY.SetBytes(yBytes)

	/*
	//if point does not lie on curve, set it to O
	if (!isValid())
		setToZero();

	return len;
	*/
	return &elem
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

// CurveElement

// TODO: Make function?

func (elem *CurveElement) IsInfinite() bool {
	return elem.DataY == nil && elem.DataY == nil
}

// satisfy Point interface
func (elem *CurveElement) X() *big.Int {
	return elem.DataX
}

func (elem *CurveElement) Y() *big.Int {
	return elem.DataY
}

func (elem *CurveElement) Copy() *CurveElement {
	newElem := *elem
	return &newElem
}

func (elem *CurveElement) Mul( n *big.Int ) *CurveElement {
	// TODO !!!
	return elem
}