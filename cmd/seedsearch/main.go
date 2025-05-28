// Package main implements the entry logic for seedsearch.
package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/chrisfenner/hashfinder/pkg/generate"
	"github.com/fatih/color"
)

var (
	margin        = flag.Int("margin", 100, "how wide to search")
	pFlag         = flag.String("p", "", "(hex) prime modulus p")
	aFlag         = flag.String("a", "", "(hex) chosen value for a (default: p-3)")
	startSeedFlag = flag.String("start_seed", "", "(hex) seed to search around for valid seeds")
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func mainErr() error {
	flag.Parse()

	if *pFlag == "" {
		return errors.New("missing value for --p")
	}
	p, err := decodeHexBig(*pFlag)
	if err != nil {
		return err
	}

	var a *big.Int
	if *aFlag == "" {
		a = big.NewInt(0).Sub(p, big.NewInt(3))
	} else {
		var err error
		a, err = decodeHexBig(*aFlag)
		if err != nil {
			return err
		}
	}

	if *startSeedFlag == "" {
		return errors.New("missing value for --start_seed")
	}
	startSeed, err := hex.DecodeString(*startSeedFlag)
	if err != nil {
		return err
	}
	return doSearch(p, a, startSeed, *margin)
}

func decodeHexBig(h string) (*big.Int, error) {
	hBytes, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(hBytes), nil
}

func doSearch(p, a *big.Int, start []byte, margin int) error {
	for i := -margin; i <= margin; i++ {
		seed := addToSeed(start, i)
		_, err := generate.PCurve(p, a, seed)
		if i == 0 {
			color.Set(color.BgBlue)
		}
		if err != nil {
			color.Set(color.FgRed)
			fmt.Printf("[BAD] %x %v", seed, err)
		} else {
			color.Set(color.FgGreen)
			fmt.Printf("[OK]  %x", seed)
		}
		color.Unset()
		fmt.Printf("\n")
	}
	return nil
}

func addToSeed(seed []byte, value int) []byte {
	result := big.NewInt(0).Add(big.NewInt(0).SetBytes(seed), big.NewInt(int64(value)))
	return result.Bytes()
}
