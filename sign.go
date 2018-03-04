package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bwesterb/go-xmssmt"

	"github.com/urfave/cli"
)

func cmdSign(c *cli.Context) error {
	var err error

	if c.NArg() != 0 {
		return cli.NewExitError("I don't expect arguments; only flags", 10)
	}

	fileInfo, err := os.Stat(c.String("privkey"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%s: no such file",
			c.String("privkey")), 11)
	}

	if fileInfo.Mode().Perm()&07 != 0 {
		return cli.NewExitError(fmt.Sprintf(
			"%s: suspicious file permission %#o",
			c.String("privkey"), fileInfo.Mode().Perm()), 12)
	}

	sk, _, lostSigs, err := xmssmt.LoadPrivateKey(c.String("privkey"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%s: %v",
			c.String("privkey"), err), 13)
	}
	defer sk.Close()

	if lostSigs != 0 {
		fmt.Fprintf(os.Stderr,
			"WARNING Lost %d XMSS[MT] signatures.\n"+
				"        This might have been caused by a crash.\n", lostSigs)
	}

	var rd io.ReadCloser
	if c.IsSet("file") {
		rd, err = os.Open(c.String("file"))

		if err != nil {
			return cli.NewExitError(fmt.Sprintf("os.Open(%s): %v",
				c.String("file"), err), 14)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Go ahead and type a message to be signed ...\n\n")
		rd = os.Stdin
	}

	sig, err := sk.SignFrom(rd)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Sign: %v", err), 15)
	}
	rd.Close()

	sigBytes, err := sig.MarshalBinary()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"Signature.MarshalBinary: %v", err), 15)
	}

	var wr io.WriteCloser
	if c.IsSet("output") || c.IsSet("file") {
		var outPath string
		if c.IsSet("output") {
			outPath = c.String("output")
		} else {
			outPath = c.String("file") + ".xmssmt-signature"
		}

		wr, err = os.Create(outPath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("os.Create(%s): %v",
				outPath, err), 16)
		}
	} else {
		wr = os.Stdout
	}

	if _, err := wr.Write(sigBytes); err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"Writing signature failed: %v", err), 16)
	}

	wr.Close()

	return nil
}
