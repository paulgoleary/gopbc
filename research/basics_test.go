package research

import (
	"fmt"
	"math/big"
	"testing"
	"crypto/sha256"
)

// val secretKeyVal = "276146606970621369032156664792541580771690346936"

func TestBasicBigness(t *testing.T) {

	secretKeyVal := "276146606970621369032156664792541580771690346936"

	qString := "8780710799663312522437781984754049815806883199414208211028653399266475630880222957078625179422662221423155858769582317459277713367317481324925129998224791"

	qVal := new(big.Int)
	qVal.SetString(qString, 0)

	testVal := new(big.Int)
	_, success := testVal.SetString(secretKeyVal, 0)
	if !success {
		t.Fail()
	}

	fmt.Println(testVal.String())

	gen := big.NewInt(3)

	kindaKey := gen.Exp(gen, testVal, qVal)

	fmt.Println(kindaKey.String())
}