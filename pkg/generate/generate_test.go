package generate

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"
)

func decodeHex(t *testing.T, h string) []byte {
	t.Helper()
	result, err := hex.DecodeString(h)
	if err != nil {
		t.Fatalf("hex.DecodeString() = %v", err)
	}
	return result
}

func decodeBigInt(t *testing.T, h string) *big.Int {
	t.Helper()
	result, err := hex.DecodeString(h)
	if err != nil {
		t.Fatalf("hex.DecodeString() = %v", err)
	}
	return big.NewInt(0).SetBytes(result)
}

func TestPCurve(t *testing.T) {
	for _, tc := range []struct {
		name       string
		seed       []byte
		p          *big.Int
		a          *big.Int
		b          *big.Int
		wantErr    error
		wantResult *Result
	}{
		{
			name: "NIST P192",
			seed: decodeHex(t, "3045AE6FC8422F64ED579528D38120EAE12196D5"),
			p:    decodeBigInt(t, "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFF"),
			a:    decodeBigInt(t, "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFC"),
			wantResult: &Result{
				R:  decodeBigInt(t, "3099D2BBBFCB2538542DCD5FB078B6EF5F3D6FE2C745DE65"),
				B1: decodeBigInt(t, "64210519E59C80E70FA7E9AB72243049FEB8DEECC146B9B1"),
			},
		},
		{
			name: "NIST P224",
			seed: decodeHex(t, "BD71344799D5C7FCDC45B59FA3B9AB8F6A948BC5"),
			p:    decodeBigInt(t, "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF000000000000000000000001"),
			a:    decodeBigInt(t, "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFE"),
			wantResult: &Result{
				R: decodeBigInt(t, "5B056C7E11DD68F40469EE7F3C7A7D74F7D121116506D031218291FB"),
				// Note that for some reason, the larger of the two solutions to r * b^2 ≡ a^3 (mod p) was chosen by NIST.
				// The other one is 4BFAF57AF3FB4C540ABECDA9AFBB4F4728402745D8F4C6BCDCAA004D.
				B1: decodeBigInt(t, "B4050A850C04B3ABF54132565044B0B7D7BFD8BA270B39432355FFB4"),
			},
		},
		{
			name: "NIST P256",
			seed: decodeHex(t, "C49D360886E704936A6678E1139D26B7819F7E90"),
			p:    decodeBigInt(t, "FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFF"),
			a:    decodeBigInt(t, "FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFC"),
			wantResult: &Result{
				R:  decodeBigInt(t, "7EFBA1662985BE9403CB055C75D4F7E0CE8D84A9C5114ABCAF3177680104FA0D"),
				B1: decodeBigInt(t, "5AC635D8AA3A93E7B3EBBD55769886BC651D06B0CC53B0F63BCE3C3E27D2604B"),
			},
		},
		{
			name: "NIST P384",
			seed: decodeHex(t, "A335926AA319A27A1D00896A6773A4827ACDAC73"),
			p:    decodeBigInt(t, "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFF0000000000000000FFFFFFFF"),
			a:    decodeBigInt(t, "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFF0000000000000000FFFFFFFC"),
			wantResult: &Result{
				R:  decodeBigInt(t, "79D1E655F868F02FFF48DCDEE14151DDB80643C1406D0CA10DFE6FC52009540A495E8042EA5F744F6E184667CC722483"),
				B1: decodeBigInt(t, "B3312FA7E23EE7E4988E056BE3F82D19181D9C6EFE8141120314088F5013875AC656398D8A2ED19D2A85C8EDD3EC2AEF"),
			},
		},
		{
			name: "NIST P521",
			seed: decodeHex(t, "D09E8800291CB85396CC6717393284AAA0DA64BA"),
			p:    decodeBigInt(t, "01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
			a:    decodeBigInt(t, "01FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC"),
			wantResult: &Result{
				R:  decodeBigInt(t, "00B48BFA5F420A34949539D2BDFC264EEEEB077688E44FBF0AD8F6D0EDB37BD6B533281000518E19F1B9FFBE0FE9ED8A3C2200B8F875E523868C70C1E5BF55BAD637"),
				B1: decodeBigInt(t, "0051953EB9618E1C9A1F929A21A0B68540EEA2DA725B99B315F3B8B489918EF109E156193951EC7E937B1652C0BD3BB1BF073573DF883D2C34F1EF451FD46B503F00"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := PCurve(tc.p, tc.a, tc.seed)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("PCurve() = %v want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Errorf("PCurve() = %v", err)
			} else {
				if !bytes.Equal(result.R.Bytes(), tc.wantResult.R.Bytes()) {
					t.Errorf("PCurve.R = %x\nwant %x", result.R.Bytes(), tc.wantResult.R.Bytes())
				}
				// Compare both B's. Just one of them has to be right.
				if !bytes.Equal(result.B1.Bytes(), tc.wantResult.B1.Bytes()) && !bytes.Equal(result.B2.Bytes(), tc.wantResult.B1.Bytes()) {
					t.Errorf("PCurve.B1 = %x\nB2 = %x\n, want %x", result.B1.Bytes(), result.B2.Bytes(), tc.wantResult.B1.Bytes())
				}
			}
		})
	}
}
