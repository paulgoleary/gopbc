package field

import "math/big"

func NAF(n *big.Int, k int) []int8 {

	// byte[] wnaf = new byte[n.bitLength() + 1];
	wnaf := make([]int8, n.BitLen() + 1)
	// short pow2wB = (short)(1 << k);
	pow2wB := 1 << (uint)(k)
	// BigInteger pow2wBI = BigInteger.valueOf((long)pow2wB);
	pow2wBI := big.NewInt((int64)(pow2wB))
	// int i = 0;

	i := 0
	length := 0
	for n.Sign() > 0 {
		if (n.Bit(0)) != 0 {
			remainder := big.Int{}
			remainder.SetBytes(n.Bytes()).Mod(n, pow2wBI) // copy n
			if remainder.Bit(k - 1) != 0 {
				wnaf[i] = (int8)(remainder.Int64() - (int64)(pow2wB))
			} else {
				wnaf[i] = (int8)(remainder.Int64())
			}

			n = n.Sub(n, big.NewInt((int64)(wnaf[i])))
			length = i
		} else {
			wnaf[i] = 0;
		}

		n.Rsh(n, 1)
		i++
	}

	length++
	wnafShort := wnaf[:length]
	return wnafShort
}