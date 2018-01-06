package crypto

import (
	"testing"
	"gopbc/pairing"
)

func TestProxyReEncryptionBasics(t *testing.T) {

	pairingParms := pairing.GetCompatParams()
	typeAPairing := pairing.MakeTypeAPairing(pairingParms)

	afgh := MakeAFGHProxyReEncryption(&typeAPairing.BasePairing)

	// get some rando secret keys
	secretKey1 := afgh.GenerateSecretKey()
	secretKey2 := afgh.GenerateSecretKey()

	// generate the matching public keys
	publicKey1 := afgh.GeneratePublicKey(secretKey1)
	publicKey2 := afgh.GeneratePublicKey(secretKey2)

	// calc the re-encryption key sk1 -> pk2
	// that is, the holder of sk1 grants the holder of sk2/pk2 a key to decrypt sk1's data
	reEncryptKey1To2 := afgh.GenerateReEncryptionKey(secretKey1, publicKey2)

	testData := []byte("Mr. Watson. Come Here. I need you.")
	testElement := afgh.makeDataElement(testData)

	// data is encrypted with the key of sk1/pk1
	c1, c2 := afgh.SecondLevelEncryption(testElement, publicKey1)

	// re-encrypted from pk1 to - essentially - pk2 by way of the re-encryption key
	c1, c2 = afgh.ReEncryption(reEncryptKey1To2, c1, c2)

	// the holder of sk2 can then completely decrypt source data
	testDecrypt := afgh.FirstLevelDecryption(secretKey2, c1, c2)

	if !testElement.IsValEqual(testDecrypt) {
		t.Errorf("Proxy re-encryption + decryption did not work: want %s, got %s", testElement.String(), testDecrypt.String())
	}
}