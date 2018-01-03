package crypto

import (
	"gopbc/pairing"
	"gopbc/field"
	"log"
)

type ProxyReEncryption struct {
	pairing.BasePairing
	Z field.PointElement
}

func (params *ProxyReEncryption) GenerateSecretKey() *field.ZElement {
	return params.Zq.NewRandomElement()
}

func (params *ProxyReEncryption) GeneratePublicKey(secretKey *field.ZElement) *field.CurveElement {
	return params.G1.GetGen().PowZn(secretKey.GetValue())
}

func (params *ProxyReEncryption) GenerateReEncryptionKey(secretKeySource *field.ZElement, publicKeyTarget *field.CurveElement) *field.CurveElement {
	// RK( a->b ) = pk_b ^(1/sk_a) = g^(b/a)
	return publicKeyTarget.PowZn(secretKeySource.Invert().GetValue())
}

// TODO: naive/POC padding ...
func (params *ProxyReEncryption) makeDataElement(data []byte) *pairing.GTFiniteElement {
	if len(data) < params.GT.LengthInBytes {
		paddingBytes := make([]byte, params.GT.LengthInBytes - len(data))
		data = append(paddingBytes, data...)
	} else if len(data) > params.GT.LengthInBytes {
		log.Panicf("Cannot make data element larger than target field allows: got %v, max %v", len(data), params.GT.LengthInBytes)
	}
	return params.GT.MakeElementFromBytes(data, params.TheMapping)
}

func (params *ProxyReEncryption) SecondLevelEncryption(data field.PointElement, publicKey field.PointElement) (field.PointElement, field.PointElement) {
	// random k \in Zq
	k := params.Zq.NewRandomElement()
	// c1 = pk_a^k
	c1 := publicKey.Pow(k.ModInt)
	// c2 = m·Z^k
	c2 := data.MulPoint(params.Z.Pow(k.ModInt))
	return c1, c2
}
