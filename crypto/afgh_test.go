package crypto

import (
	"testing"
	"gopbc/pairing"
	"gopbc/field"
)

func TestProxyReEncryptionBasics(t *testing.T) {

	pairingParms := pairing.GetCompatParams()
	typeAPairing := pairing.MakeTypeAPairing(pairingParms)

	// TODO: obv move someplace else ... :/
	afgh := new(ProxyReEncryption)
	afgh.BasePairing = typeAPairing.BasePairing
	afgh.Z = afgh.TheMapping.Pairing(afgh.G1.GetGen(), afgh.G1.GetGen())

	secretKey1 := afgh.GenerateSecretKey()
	secretKey2 := afgh.GenerateSecretKey()
	field.Trace(secretKey1, secretKey2)

	publicKey1 := afgh.GeneratePublicKey(secretKey1)
	publicKey2 := afgh.GeneratePublicKey(secretKey2)
	field.Trace(publicKey1, publicKey2)

	reEncryptKey1To2 := afgh.GenerateReEncryptionKey(secretKey1, publicKey2)
	field.Trace(reEncryptKey1To2)

	testData := []byte("Mr. Watson. Come Here. I need you.")
	testElement := afgh.makeDataElement(testData)
	field.Trace(testElement)

	c1, c2 := afgh.SecondLevelEncryption(testElement, publicKey1)
	field.Trace(c1, c2)
}