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
	DataX     *BigInt // TODO: perhaps X and Y should be elements of the target field, as in PBC/JPBC?
	DataY     *BigInt
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

	elem.DataX = new(BigInt)
	elem.DataX.SetBytes(xBytes)

	elem.DataY = new(BigInt)
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

func (elem *CurveElement) IsInf() bool {
	return elem.DataY == nil && elem.DataY == nil
}

func (elem *CurveElement) SetInf() {
	elem.DataX = nil
	elem.DataY = nil
}

// satisfy PointLike interface
func (elem *CurveElement) X() *BigInt {
	return elem.DataX
}

func (elem *CurveElement) Y() *BigInt {
	return elem.DataY
}

// TODO: need to fully copy X and Y
func (elem *CurveElement) Copy() *CurveElement {
	newElem := *elem
	return &newElem
}

func (elem *CurveElement) Mul( n *big.Int ) *CurveElement {
	// TODO !!!
	return elem
}

func (elem *CurveElement) set( in *CurveElement ) {
	elem.DataX = in.DataX
	elem.DataY = in.DataY
}

// TODO: wrapper for big.Int?
type BigInt big.Int

func MakeBigInt(x int64) *BigInt {
	return (*BigInt)(big.NewInt(x))
}

func MakeBigIntStr(x string) *BigInt {
	ret := big.Int{}
	ret.SetString(x, 10)
	return (*BigInt)(&ret)
}

func (bi *BigInt) isZero() bool {
	return (*big.Int)(bi).Cmp(ZERO) == 0
}

func (bi *BigInt) copy() *BigInt {
	newBigInt := new(BigInt)
	(*big.Int)(newBigInt).SetBytes((*big.Int)(bi).Bytes())
	return newBigInt
}

func (bi *BigInt) SetBytes(bytes []byte) {
	(*big.Int)(bi).SetBytes(bytes)
}

func (bi *BigInt) IsEqual(in *BigInt) bool {
	return (*big.Int)(bi).Cmp((*big.Int)(in)) == 0
}

func (bi *BigInt) sub(in *BigInt) *BigInt {
	(*big.Int)(bi).Sub((*big.Int)(bi), (*big.Int)(in))
	return bi
}

func (bi *BigInt) mul(in *BigInt) *BigInt {
	(*big.Int)(bi).Mul((*big.Int)(bi), (*big.Int)(in))
	return bi
}

func (bi *BigInt) square() *BigInt {
	(*big.Int)(bi).Mul((*big.Int)(bi), (*big.Int)(bi))
	return bi
}

func (bi *BigInt) invert(mod *big.Int) *BigInt {
	(*big.Int)(bi).ModInverse((*big.Int)(bi), mod)
	return bi
}

func (bi *BigInt) String() string {
	return (*big.Int)(bi).String()
}

func (elem *CurveElement) twiceInternal() *CurveElement {
	// TODO !!!
	return elem
}

func (elem *CurveElement) mul(elemIn *CurveElement) *CurveElement {

	if elem.IsInf() {
		elem.set(elemIn)
		return elem
	}

	if elemIn.IsInf() {
		return elem
	}

	if elem.DataX.IsEqual(elemIn.DataX) {
		if elem.DataY.IsEqual(elemIn.DataY) {
			if elem.DataY.IsEqual(BI_ZERO) {
				elem.SetInf()
				return elem
			} else {
				elem.twiceInternal()
				return elem
			}
		}
		elem.SetInf()
		return elem
	}

	// P1 != P2, so the slope of the line L through P1 and P2 is
	// lambda = (y2-y1)/(x2-x1)
	// Element lambda = element.y.duplicate().sub(y).mul(element.x.duplicate().sub(x).invert());
	lambdaNumer := elemIn.DataY.copy().sub(elem.DataY)
	lambdaDenom := elemIn.DataX.copy().sub(elem.DataX)
	lambda := lambdaNumer.mul(lambdaDenom.invert(elem.ElemField.order))

	// x3 = lambda^2 - x1 - x2
	// Element x3 = lambda.duplicate().square().sub(x).sub(element.x);
	x3 := lambda.copy().square().sub(elem.DataX).sub(elemIn.DataX)

	//y3 = (x1-x3)lambda - y1
	// Element y3 = x.duplicate().sub(x3).mul(lambda).sub(y);
	y3 := elem.DataX.copy().sub(x3).mul(lambda).sub(elem.DataY)

	// x.set(x3);
	// y.set(y3);
	elem.DataX = x3
	elem.DataY = y3

	return elem
}
