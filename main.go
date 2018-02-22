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
		},
		{
			Name:   "speed",
			Usage:  "Benchmark XMSS[MT] instances",
			Action: cmdSpeed,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "benchmark instance named `NAME`",
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
