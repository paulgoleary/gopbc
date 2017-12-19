package pairing

import "gobdc/field"

type JacobPoint struct {
	field.PointLike
	z *field.BigInt
}
