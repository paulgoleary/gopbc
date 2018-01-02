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