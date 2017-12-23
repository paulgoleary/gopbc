package pairing

import (
	"gobdc/field"
	"fmt"
)

type JacobPoint struct {
	x *field.BigInt
	y *field.BigInt
	z *field.BigInt
}

func (jp *JacobPoint) String() string {
	return fmt.Sprintf("[%s,\n%s,\n%s]", jp.x.String(), jp.y.String(), jp.z.String())
}

func (jp *JacobPoint) freeze() {
	jp.x.Freeze()
	jp.y.Freeze()
	jp.z.Freeze()
}
