package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bwesterb/go-xmssmt"

	"github.com/urfave/cli"
)

func cmdGenerate(c *cli.Context) error {
	var err error
	params, err := xmssmt.ParamsFromName2(c.String("alg"))

	if err != nil {
		return cli.NewExitError(fmt.Sprintf(
			"There is no XMSS[MT] instance %s: %v", c.String("alg"), err), 1)
	}

	if c.NArg() != 0 {
		return cli.NewExitError("I don't expect arguments; only flags", 10)
	}

	if c.IsSet("n") {
		params.N = uint32(c.Int("n"))
	}

	if c.IsSet("w") {
		params.WotsW = uint16(c.Int("w"))
	}

	if c.IsSet("d") {
		params.D = uint32(c.Int("d"))
	}

	if c.IsSet("full-height") {
		params.FullHeight = uint32(c.Int("full-height"))
	}

	if c.IsSet("hash") {
		switch c.String("hash") {
		case "sha2":
			params.Func = xmssmt.SHA2
		case "shake":
			params.Func = xmssmt.SHAKE
		default:
			return cli.NewExitError(fmt.Sprintf(
				"The hash function %s is not supported", c.String("hash")), 2)
		}
	}

	ctx, err := xmssmt.NewContext(*params)

	if err != nil {
		return cli.NewExitError(err, 3)
	}

	if !c.Bool("force") {
		if _, err = os.Stat(c.String("privkey")); !os.IsNotExist(err) {
			return cli.NewExitError(fmt.Sprintf(
				"%s: already exists", c.String("privkey")), 8)
		}
		if _, err = os.Stat(c.String("pubkey")); !os.IsNotExist(err) {
			return cli.NewExitError(fmt.Sprintf(
				"%s: already exists", c.String("pubkey")), 9)
		}
	}

	sk, pk, err := ctx.GenerateKeyPair(c.String("privkey"))

	if err != nil {
		return cli.NewExitError(err, 4)
	}

	err = sk.Close()
	if err != nil {
		return cli.NewExitError(err, 5)
	}

	pkBytes, err := pk.MarshalBinary()
	if err != nil {
		return cli.NewExitError(err, 6)
	}
	pkBase64, _ := pk.MarshalText()

	err = ioutil.WriteFile(c.String("pubkey"), pkBytes, 0644)
	if err != nil {
		return cli.NewExitError(err, 7)
	}

	fmt.Printf("Successfully generated keypair.\n\n")
	fmt.Printf("Public key: %s\n", pkBase64)

	return nil
}
