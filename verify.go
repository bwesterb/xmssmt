package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bwesterb/go-xmssmt"

	"github.com/urfave/cli"
)

func cmdVerify(c *cli.Context) error {
	var err error
	var pk xmssmt.PublicKey
	var sig xmssmt.Signature

	if c.NArg() != 0 {
		return cli.NewExitError("I don't expect arguments; only flags", 10)
	}

	pkBytes, err := ioutil.ReadFile(c.String("pubkey"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"os.Open(%s): %v", c.String("pubkey"), err), 17)
	}

	if err := pk.UnmarshalBinary(pkBytes); err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"%s: %v", c.String("pubkey"), err), 17)
	}

	var sigPath string
	if c.IsSet("signature") {
		sigPath = c.String("signature")
	} else if c.IsSet("file") {
		sigPath = c.String("file") + ".xmssmt-signature"
	} else {
		return cli.NewExitError(
			"Either --file or --signature should be provided", 18)
	}

	sigBytes, err := ioutil.ReadFile(sigPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"os.Open(%s): %v", sigPath, err), 19)
	}
	if err := sig.UnmarshalBinary(sigBytes); err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"%s: %v", sigPath, err), 19)
	}

	var rd io.ReadCloser
	if c.IsSet("file") {
		rd, err = os.Open(c.String("file"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("os.Open(%s): %v",
				c.String("file"), err), 20)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Go ahead and type the message to be verified ...\n\n")
		rd = os.Stdin
	}

	valid, err := pk.VerifyFrom(&sig, rd)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Verify: %v", err), 21)
	}

	if !valid {
		return cli.NewExitError(fmt.Sprintf("Signature is *not* valid: %v", err), 22)
	}

	fmt.Fprintf(os.Stderr, "Signature is valid\n")

	rd.Close()

	return nil
}
