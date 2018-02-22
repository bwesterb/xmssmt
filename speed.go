package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func benchmarkSign(dir string, b *testing.B) {
	sk, _, _, err := xmssmt.LoadPrivateKey(dir + "/key")
	if err != nil {
		b.Fatalf("LoadPrivateKey: %v", err)
	}
	defer sk.Close()
	sk.BorrowExactly(1000000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sk.DangerousSetSeqNo(xmssmt.SignatureSeqNo(0))
		sk.Sign([]byte("test message"))
	}
}

func benchmarkVerify(dir string, b *testing.B) {
	sk, pk, _, err := xmssmt.LoadPrivateKey(dir + "/key")
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
	nameFlag := c.String("name")
	if nameFlag != "" {
		toTest = []string{nameFlag}
	} else {
		toTest = xmssmt.ListNames()
	}

	for _, name := range toTest {
		fmt.Printf("%s\n", name)
		params := xmssmt.ParamsFromName(name)
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
		res = testing.Benchmark(func(b *testing.B) {
			benchmarkSign(dir, b)
		})
		fmt.Printf(" sign:   %20s %s\n",
			humanize.SI(float64(res.NsPerOp())/1e9, "s"),
			res.MemString())
		res = testing.Benchmark(func(b *testing.B) {
			benchmarkVerify(dir, b)
		})
		fmt.Printf(" verify: %20s %s\n\n",
			humanize.SI(float64(res.NsPerOp())/1e9, "s"),
			res.MemString())
	}

	return nil
}
