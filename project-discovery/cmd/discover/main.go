package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Ticket ticketCmd `cmd:"" help:"ticket operations"`
	Scan   scanCmd   `cmd:"" help:"repo facts scan"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("discover"),
		kong.Description("project-discovery CLI"),
		kong.UsageOnError(),
	)
	ctx.FatalIfErrorf(ctx.Run())
}
