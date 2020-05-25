package main

import (
	"fmt"
	"os"

	"github.com/bwesterb/go-xmssmt"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func cmdAlgs(c *cli.Context) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"name",
		"oid",
		"#sigs",
		"sigSize",
		"cache size"})
	names := xmssmt.ListNames()
	if c.Bool("non-std") {
		names = xmssmt.ListNames2()
	}
	for _, name := range names {
		ctx, _ := xmssmt.NewContextFromName2(name)
		var cacheSize uint64
		if ctx.MT() {
			cacheSize = uint64((ctx.Params().D+1)*ctx.Params().N) * (1 << uint64(ctx.Params().FullHeight/ctx.Params().D))
		} else {
			cacheSize = (1 << uint64(ctx.Params().FullHeight)) *
				uint64(ctx.Params().N)
		}
		table.Append([]string{
			name,
			fmt.Sprintf("%d", ctx.Oid()),
			fmt.Sprintf("2^%d", ctx.Params().FullHeight),
			humanize.Bytes(uint64(ctx.SignatureSize())),
			humanize.Bytes(uint64(cacheSize)),
		})
	}
	table.Render()

	return nil
}
