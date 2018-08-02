package pairing

import (
	"github.com/paulgoleary/gopbc/field"
	"fmt"
)

type JacobPoint struct {
	x *field.ModInt
	y *field.ModInt
	z *field.ModInt
}

func (jp *JacobPoint) String() string {
	return fmt.Sprintf("[%s,\n%s,\n%s]", jp.x.String(), jp.y.String(), jp.z.String())
}

func (jp *JacobPoint) freeze() {
	jp.x.Freeze()
	jp.y.Freeze()
	jp.z.Freeze()
}
