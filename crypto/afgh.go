package crypto

import (
	"github.com/paulgoleary/gopbc/pairing"
	"github.com/paulgoleary/gopbc/field"
	"log"
)

type ProxyReEncryption struct {
	pairing.BasePairing
	Z field.PointElement
}

func MakeAFGHProxyReEncryption(thePairing *pairing.BasePairing) *ProxyReEncryption {
	// TODO: obv move someplace else ... :/
	afgh := new(ProxyReEncryption)
	afgh.BasePairing = *thePairing
	afgh.Z = afgh.TheMapping.Pairing(afgh.G1.GetGen(), afgh.G1.GetGen())
	return afgh
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

func (params *ProxyReEncryption) ReEncryption(reEncryptKey1To2 field.PointElement, c1 field.PointElement, c2 field.PointElement) (field.PointElement, field.PointElement) {
	// c1 = ( e(c1, rk) = e(pk_a^k, (pk_b ^(1/sk_a) = g^(b/a)))
	c1 = params.TheMapping.Pairing(c1, reEncryptKey1To2)
	return c1, c2
}

func (params *ProxyReEncryption) FirstLevelDecryption(targetSecretKey *field.ZElement, c1 field.PointElement, c2 field.PointElement ) field.PointElement {
	// c1 = ( e(c1, rk) = e(pk_a^k, g^(b/a)) = e(g^ak, g^(b/a)) = e(g,g)^ak(b/a) = Z^bk
	// c2 = m·Z^k
	// c1x = Z^bk ^ 1/b = Z^k
	c1x := c1.Pow(targetSecretKey.Invert())
	// m·Z^k * (1/Z^k) = m
	return c2.MulPoint(c1x.Invert())
}