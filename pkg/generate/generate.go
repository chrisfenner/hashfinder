// Package generate implements the A.3 curve generation algorithm from https://safecurves.cr.yp.to/grouper.ieee.org/groups/1363/private/x9-62-09-20-98.pdf.
package generate

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrBadSeed = errors.New("bad seed")
)

type Result struct {
	R      *big.Int
	B1, B2 *big.Int
}

// PCurve computes NIST P-curve parameters for a given a, p, and SEED.
func PCurve(p, a *big.Int, seed []byte) (*Result, error) {
	sCount := (p.BitLen() - 1) / 160
	hBits := p.BitLen() - 160*sCount
	hBytes := (hBits + 7) / 8
	var rBuf bytes.Buffer

	// Compute W_0 which is a truncated hash of SEED.
	w0Buf := sha1.Sum(seed)
	w0 := w0Buf[20-hBytes:]
	for i := 0; i <= (160-hBits)%8; i++ {
		w0[0] &= ^(1 << (7 - i))
	}
	rBuf.Write(w0)

	// Compute the rest of the W's.
	for i := 1; i <= sCount; i++ {
		seedInt := big.NewInt(0).Add(big.NewInt(0).SetBytes(seed), big.NewInt(int64(i)))
		w := sha1.Sum(seedInt.Bytes())
		rBuf.Write(w[:])
	}
	r := big.NewInt(0).SetBytes(rBuf.Bytes())

	// Compute B such that r * b^2 ≡ a^3 (mod p)
	a3 := big.NewInt(0).Exp(a, big.NewInt(3), p)
	if a3 == nil {
		return nil, fmt.Errorf("%w: error computing a^3", ErrBadSeed)
	}
	rinv := big.NewInt(0).ModInverse(r, p)
	if rinv == nil {
		return nil, fmt.Errorf("%w: error computing r-inverse", ErrBadSeed)
	}
	b2 := big.NewInt(0).Mod(big.NewInt(0).Mul(a3, rinv), p)
	if b2 == nil {
		return nil, fmt.Errorf("%w: error computing b^2", ErrBadSeed)
	}

	// Offer both square roots.
	b := big.NewInt(0).ModSqrt(b2, p)
	if b == nil {
		return nil, fmt.Errorf("%w: error computing b", ErrBadSeed)
	}
	bneg := big.NewInt(0).Sub(p, b)

	// Check if 4a^3 + 27b^2 ≡ 0 (mod p)
	check := big.NewInt(0).Mod(big.NewInt(0).Add(big.NewInt(0).Mul(a3, big.NewInt(4)), big.NewInt(0).Mul(b2, big.NewInt(27))), p)
	if check.IsInt64() && check.Int64() == 0 {
		return nil, fmt.Errorf("%w: 4a^3 + 27b^2 ≡ 0 (mod p)", ErrBadSeed)
	}

	return &Result{
		R:  r,
		B1: b,
		B2: bneg,
	}, nil
}
