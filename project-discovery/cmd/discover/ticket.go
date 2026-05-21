package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/urso/claudev/project-discovery/pkg/ticket"
)

type ticketCmd struct {
	NextID      nextIDCmd      `cmd:"" name:"next-id" help:"print next ticket id"`
	New         newCmd         `cmd:"" help:"create a new ticket"`
	List        listCmd        `cmd:"" help:"list tickets"`
	Get         getCmd         `cmd:"" help:"get a ticket by id"`
	Update      updateCmd      `cmd:"" help:"update a ticket"`
	Search      searchCmd      `cmd:"" help:"search tickets"`
	FindOverlap findOverlapCmd `cmd:"" name:"find-overlap" help:"detect overlapping tickets"`
	Recall      recallCmd      `cmd:"" help:"recall-first search for existing knowledge"`
	Status      statusCmd      `cmd:"" help:"show ticket status summary"`
	Next        nextCmd        `cmd:"" help:"suggest next ticket to work on"`
	Tags        tagsCmd        `cmd:"" help:"list tags"`
}

func tickets() (ticket.TicketFolder, error) {
	d, err := ticket.FindDiscovery(".")
	if err != nil {
		return ticket.TicketFolder{}, err
	}
	return d.Tickets(), nil
}

type nextIDCmd struct{}

func (nextIDCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	id, err := ticket.NextID(t)
	if err != nil {
		return err
	}
	fmt.Println(id)
	return nil
}

type newCmd struct {
	Title     string   `required:"" help:"ticket title"`
	Slug      string   `help:"explicit slug ([a-z0-9-]); derived from title if empty"`
	Scope     string   `help:"one-line scope"`
	Tag       []string `help:"tag (repeatable)"`
	Parent    []string `help:"parent ticket id (repeatable)"`
	Intention string   `help:"why this ticket exists"`
	Body      string   `help:"body content; use '-' to read from stdin"`
}

func (c newCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	var body []byte
	switch c.Body {
	case "":
	case "-":
		body, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
	default:
		body = []byte(c.Body)
	}
	doc, err := ticket.New(t, ticket.NewParams{
		Title:     c.Title,
		Slug:      c.Slug,
		Body:      body,
		Scope:     c.Scope,
		Tags:      c.Tag,
		Parents:   c.Parent,
		Intention: c.Intention,
	})
	if err != nil {
		return err
	}
	if err := t.Write(doc); err != nil {
		return err
	}
	fmt.Printf("%s\t%s\n", doc.Frontmatter.ID, doc.Path)
	return nil
}

type listCmd struct {
	Status string `help:"filter by status"`
	Tag    string `help:"filter by tag"`
}

func (c listCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	rows, err := ticket.List(t, ticket.ListFilter{Status: c.Status, Tag: c.Tag})
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", r.ID, r.Status, r.Title)
	}
	return tw.Flush()
}

type getCmd struct {
	ID string `arg:"" help:"ticket id (e.g. t-0001 or 1)"`
}

func (c getCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	id, err := ticket.ParseID(c.ID)
	if err != nil {
		return err
	}
	doc, err := t.Read(id)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(doc.Path)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}

type updateCmd struct {
	ID     string   `arg:"" help:"ticket id"`
	Status string   `help:"new status (open|active|done|dropped|blocked)"`
	Log    []string `help:"log entry (repeatable)"`
}

func (c updateCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	id, err := ticket.ParseID(c.ID)
	if err != nil {
		return err
	}
	if c.Status == "" && len(c.Log) == 0 {
		return fmt.Errorf("nothing to update: pass --status and/or --log")
	}
	if c.Status != "" {
		if err := ticket.UpdateStatus(t, id, c.Status); err != nil {
			return err
		}
	}
	if len(c.Log) > 0 {
		if err := ticket.AppendLogs(t, id, c.Log); err != nil {
			return err
		}
	}
	return nil
}

type searchCmd struct {
	Query string   `arg:"" help:"search query"`
	Tag   []string `help:"require tag (repeatable)"`
	Scope string   `help:"scope filter"`
}

func (c searchCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	hits, err := ticket.Search(t, c.Query, ticket.SearchFilter{Tags: c.Tag, Scope: c.Scope})
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, h := range hits {
		fmt.Fprintf(tw, "%s\t%.4f\t%s\t%s\n", h.ID, h.Score, h.Status, h.Title)
	}
	return tw.Flush()
}

type findOverlapCmd struct {
	Intent string   `required:"" help:"intent of the new ticket"`
	Scope  string   `help:"scope of the new ticket"`
	Tag    []string `help:"tag of the new ticket (repeatable)"`
}

func (c findOverlapCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	hits, err := ticket.FindOverlap(t, ticket.OverlapParams{
		Intent: c.Intent,
		Scope:  c.Scope,
		Tags:   c.Tag,
	})
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, h := range hits {
		fmt.Fprintf(tw, "%s\t%.6f\t%s\t%s\n", h.ID, h.Score, h.Status, h.Title)
	}
	return tw.Flush()
}

type recallCmd struct {
	Intent string   `arg:"" help:"intent to search for"`
	Scope  string   `help:"scope filter"`
	Tag    []string `help:"tag filter (repeatable)"`
}

func (c recallCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	result, err := ticket.Recall(t, c.Intent, c.Scope, c.Tag)
	if err != nil {
		return err
	}
	fmt.Println(result.Type)
	if len(result.Matches) > 0 {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, h := range result.Matches {
			fmt.Fprintf(tw, "%s\t%.4f\t%s\t%s\n", h.ID, h.Score, h.Status, h.Title)
		}
		tw.Flush()
	}
	return nil
}

type statusCmd struct{}

func (statusCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	s, err := ticket.Status(t)
	if err != nil {
		return err
	}
	fmt.Print(ticket.FormatStatus(s))
	return nil
}

type nextCmd struct {
	N int `default:"3" help:"number of suggestions"`
}

func (c nextCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	suggestions, err := ticket.Next(t, c.N)
	if err != nil {
		return err
	}
	if len(suggestions) == 0 {
		fmt.Println("No open or active tickets.")
		return nil
	}
	for i, s := range suggestions {
		fmt.Printf("%d. %s: %s — %s\n", i+1, s.ID, s.Title, s.Reason)
	}
	return nil
}

type tagsCmd struct{}

func (tagsCmd) Run() error {
	t, err := tickets()
	if err != nil {
		return err
	}
	rows, err := ticket.Tags(t)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintf(tw, "%s\t%d\n", r.Tag, r.Count)
	}
	return tw.Flush()
}
