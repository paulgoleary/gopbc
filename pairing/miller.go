package pairing

import (
	"gobdc/field"
	"fmt"
)

type JacobPoint struct {
	field.PointLike
	z *field.BigInt
}

func (jp *JacobPoint) String() string {
	return fmt.Sprintf("[%s,\n%s,\n%s]", jp.DataX.String(), jp.DataY.String(), jp.z.String())
}

func (jp *JacobPoint) freeze() {
	jp.DataX.Freeze()
	jp.DataY.Freeze()
	jp.z.Freeze()
}
