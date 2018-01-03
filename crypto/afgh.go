package crypto

import (
	"gopbc/pairing"
	"gopbc/field"
)

type ProxyReEncryption struct {
	pairing.BasePairing
}

func (params *ProxyReEncryption) GenerateSecretKey() *field.ZrElement {
	randInt := field.GetRandomInt(params.Zr.FieldOrder)
	return params.Zr.NewElement(randInt)
}

func (params *ProxyReEncryption) GeneratePublicKey(secretKey *field.ZrElement) *field.CurveElement {
	return params.G1.GetGen().PowZn(secretKey.GetValue())
}

func (params *ProxyReEncryption) GenerateReEncryptionKey(secretKeySource *field.ZrElement, publicKeyTarget *field.CurveElement) *field.CurveElement {
	// RK( a->b ) = pk_b ^(1/sk_a) = g^(b/a)
	return publicKeyTarget.PowZn(secretKeySource.Invert().GetValue())
}
