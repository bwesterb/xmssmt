package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/urfave/cli"
)

func main() {
	var cpuProfileFile *os.File

	app := cli.NewApp()

	app.Version = "1.1.0"
	app.Usage = "Create and verify XMSS[MT] signatures"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cpuprofile, p",
			Usage: "write cpu profile to `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "algs",
			Usage:  "List XMSS[MT] instances",
			Action: cmdAlgs,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "non-std, n",
					Usage: "Include instances which are not listed in the RFC " +
						"or NIST standard",
				},
			},
		},
		{
			Name:   "generate",
			Usage:  "Generate an XMSS[MT] keypair",
			Action: cmdGenerate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alg, a",
					Usage: "XMSS[MT] named instance to use, see `xmssmt algs'",
					Value: "XMSSMT-SHAKE_40/4_256",
				},
				cli.IntFlag{
					Name:  "n",
					Usage: "Override security parameter n",
					Value: 32,
				},
				cli.IntFlag{
					Name:  "w",
					Usage: "Override Winternitz parameter w",
					Value: 16,
				},
				cli.IntFlag{
					Name:  "full-height, t",
					Usage: "Override full tree height parameter",
					Value: 40,
				},
				cli.IntFlag{
					Name:  "d",
					Usage: "Override height-of-hypertree paramater d",
					Value: 4,
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Override existing files",
				},
				cli.StringFlag{
					Name:  "hash, H",
					Usage: "Override hash function to use. (shake, shake256 or sha2)",
					Value: "shake",
				},
				cli.StringFlag{
					Name:  "prf, P",
					Usage: "Override prf function to use. (rfc or nist)",
				},
				cli.StringFlag{
					Name:  "privkey, s",
					Usage: "Path to store private key at",
					Value: "xmssmt.key",
				},
				cli.StringFlag{
					Name:  "pubkey, p",
					Usage: "Path to store public key at",
					Value: "xmssmt.pub",
				},
			},
		},
		{
			Name:   "sign",
			Usage:  "Create an XMSS[MT] signature",
			Action: cmdSign,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "privkey, s",
					Usage: "Use private key stored at `FILE`",
					Value: "xmssmt.key",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Create a signature of `FILE`",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Write signature to `FILE`",
				},
			},
		},
		{
			Name:   "verify",
			Usage:  "Verifies an XMSS[MT] signature",
			Action: cmdVerify,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "pubkey, p",
					Usage: "Path to read public key from `FILE`",
					Value: "xmssmt.pub",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Reads message from `FILE`",
				},
				cli.StringFlag{
					Name:  "signature, S",
					Usage: "Reads signature from `FILE`",
				},
			},
		},
		{
			Name:   "speed",
			Usage:  "Benchmark XMSS[MT] instances",
			Action: cmdSpeed,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alg, a",
					Usage: "Benchmark instance named `NAME`",
				},
				cli.BoolFlag{
					Name:  "cwd",
					Usage: "Look for existing key in current working directory",
				},
				cli.BoolFlag{
					Name: "non-std, n",
					Usage: "Include instances which are not listed in the " +
						"RFC or NIST standard",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		if cpuProfilePath := c.String("cpuprofile"); cpuProfilePath != "" {
			var err error
			cpuProfileFile, err = os.Create(cpuProfilePath)
			if err != nil {
				log.Fatalf("os.Create(): %v", err)
			}
			pprof.StartCPUProfile(cpuProfileFile)
		}
		return nil
	}

	app.After = func(c *cli.Context) error {
		if cpuProfileFile != nil {
			pprof.StopCPUProfile()
			cpuProfileFile.Close()
		}
		return nil
	}

	app.Run(os.Args)
}
