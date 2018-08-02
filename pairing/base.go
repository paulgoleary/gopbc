package pairing

import (
	"encoding/base64"
	"fmt"
	"github.com/paulgoleary/gopbc/field"
	"math/big"
	"strconv"
)

const (
	aType = "a"
)

type PairingParameters map[string]string

// TODO
type PreProcessing interface {}

// TODO!
type BasePairing struct {
	// protected SecureRandom random;
	G1, G2     *field.CurveField
	GT         *GTFiniteField
	Zq         *field.ZField
	TheMapping Mapping
}

func (params PairingParameters) getInt(paramName string) int {
	ret, err := strconv.Atoi(params[paramName])
	if err != nil {
		panic(err.Error()) // TODO: additional error context?
	}
	return ret
}

// TODO: error handling
func (params PairingParameters) getBigInt(paramName string) *big.Int {
	paramVal := params[paramName]
	ret := new(big.Int)
	ret.SetString(paramVal, 10)
	return ret
}

func (params PairingParameters) getBytes(paramName string, defaultVal []byte) []byte {
	paramVal := params[paramName]
	if paramVal == "" {
		return defaultVal
	}
	decoded, err := base64.StdEncoding.DecodeString(paramVal)
	if err != nil {
		panic(fmt.Sprintf("Could not base64 decode param name '%s', value: %s", paramName, paramVal))
	}
	return decoded
}

func (params PairingParameters) getString(paramName string, defaultVal string) string {
	paramVal := params[paramName]
	if paramVal == "" {
		return defaultVal
	}
	return paramVal
}
