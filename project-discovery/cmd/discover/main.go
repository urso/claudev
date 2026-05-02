package main

import (
	"errors"

	"github.com/alecthomas/kong"
)

type notImplCmd struct{}

func (notImplCmd) Run() error { return errors.New("not implemented") }

type ticketCmd struct {
	NextID      notImplCmd `cmd:"" name:"next-id" help:"print next ticket id"`
	New         notImplCmd `cmd:"" help:"create a new ticket"`
	List        notImplCmd `cmd:"" help:"list tickets"`
	Get         notImplCmd `cmd:"" help:"get a ticket by id"`
	Update      notImplCmd `cmd:"" help:"update a ticket"`
	Search      notImplCmd `cmd:"" help:"search tickets"`
	FindOverlap notImplCmd `cmd:"" name:"find-overlap" help:"detect overlapping tickets"`
	Tags        notImplCmd `cmd:"" help:"list tags"`
}

type scanCmd struct {
	Repo  notImplCmd `cmd:"" help:"scan repo facts"`
	Churn notImplCmd `cmd:"" help:"scan file churn"`
	Deps  notImplCmd `cmd:"" help:"scan dependencies"`
}

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
