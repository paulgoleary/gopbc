package crypto

import (
	"gopbc/pairing"
	"gopbc/field"
	"log"
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
