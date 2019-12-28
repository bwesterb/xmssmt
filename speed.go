package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/bwesterb/go-xmssmt"

	"github.com/dustin/go-humanize"
	"github.com/urfave/cli"
)

func benchmarkKeygen(dir string, params xmssmt.Params, b *testing.B) {
	ctx, _ := xmssmt.NewContext(params)

	pubSeed := make([]byte, params.N)
	skSeed := make([]byte, params.N)
	skPrf := make([]byte, params.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sk, _, _ := ctx.Derive(dir+"/key", pubSeed, skSeed, skPrf)
		sk.Close()
	}
}

func benchmarkSign(path string, b *testing.B) {
	sk, _, _, err := xmssmt.LoadPrivateKey(path)
	if err != nil {
		b.Fatalf("LoadPrivateKey: %v", err)
	}
	defer sk.Close()
	sk.BorrowExactly(1000000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sk.DangerousSetSeqNo(xmssmt.SignatureSeqNo(0))
		sk.Sign([]byte(strconv.Itoa(i)))
	}
}

func benchmarkVerify(path string, b *testing.B) {
	sk, pk, _, err := xmssmt.LoadPrivateKey(path)
	if err != nil {
		b.Fatalf("LoadPrivateKey: %v", err)
	}
	defer sk.Close()
	sk.BorrowExactly(1000000)
	sk.DangerousSetSeqNo(xmssmt.SignatureSeqNo(0))
	sig, _ := sk.Sign([]byte("test message"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pk.Verify(sig, []byte("test message"))
	}
}

func cmdSpeed(c *cli.Context) error {
	var toTest []string
	var path string

	algFlag := c.String("alg")
	if algFlag != "" {
		toTest = []string{algFlag}
	} else {
		toTest = xmssmt.ListNames()
		if c.Bool("non-rfc") {
			toTest = xmssmt.ListNames2()
		}
	}

	for _, name := range toTest {
		fmt.Printf("%s\n", name)
		params, err := xmssmt.ParamsFromName2(name)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf(
				"There is no XMSS[MT] instance %s: %v", name, err), 1)
		}
		if c.Bool("cwd") {
			path = name
			if _, err := os.Stat(path); os.IsNotExist(err) {
				ctx, _ := xmssmt.NewContext(*params)
				sk, _, _ := ctx.GenerateKeyPair(path)
				sk.Close()
			}
		} else {
			dir, err := ioutil.TempDir("", "go-xmssmt-tests")
			if err != nil {
				log.Fatalf("TempDir: %v", err)
			}
			defer os.RemoveAll(dir)
			res := testing.Benchmark(func(b *testing.B) {
				benchmarkKeygen(dir, *params, b)
			})
			fmt.Printf(" keygen: %20s %s\n",
				humanize.SI(float64(res.NsPerOp())/1e9, "s"),
				res.MemString())
			path = dir + "/key"
		}
		res := testing.Benchmark(func(b *testing.B) {
			benchmarkSign(path, b)
		})
		fmt.Printf(" sign:   %20s %s\n",
			humanize.SI(float64(res.NsPerOp())/1e9, "s"),
			res.MemString())
		res = testing.Benchmark(func(b *testing.B) {
			benchmarkVerify(path, b)
		})
		fmt.Printf(" verify: %20s %s\n\n",
			humanize.SI(float64(res.NsPerOp())/1e9, "s"),
			res.MemString())
	}

	return nil
}
