package crypto

import (
	"testing"
	"github.com/paulgoleary/gopbc/pairing"
	"crypto/sha256"
)

func hashMessage(message []byte) []byte {
	h := sha256.New()
	h.Write(message)
	return h.Sum(nil)
}

// only a quick test to demonstrate the basics of BLS style short signatures
// https://en.wikipedia.org/wiki/Boneh%E2%80%93Lynn%E2%80%93Shacham
func TestBLSSignatureBasics(t *testing.T) {

	pairingParms := pairing.GetCompatParams()
	typeAPairing := pairing.MakeTypeAPairing(pairingParms)

	// get rando secret key
	secretKey := typeAPairing.Zq.NewRandomElement()
	// corresponding public key ...
	publicKey := typeAPairing.G1.GetGen().PowZn(secretKey.GetValue())

	testData := []byte("Mr. Watson. Come Here. I need you.")

	h := hashMessage(testData)
	elemHash := typeAPairing.G1.MakeElementFromHash(h)

	// sig is just h^sk - secret key of the signer
	sig := elemHash.Pow(secretKey.ModInt)

	// we'll pretend that the verifier hashed the message as well

	// e(sig, g) = e(h^sk, g) = e(h, g)^sk
	p1 := typeAPairing.TheMapping.Pairing(sig, typeAPairing.G1.GetGen())
	// e(h, pk) = e(h, g^sk) = e(h, g)^sk - public key of the signer
	p2 := typeAPairing.TheMapping.Pairing(elemHash, publicKey)

	if !p1.IsValEqual(p2) {
		t.Errorf("Signature pairings are not equal: want %s, got %s", p1.String(), p2.String() )
	}
}
