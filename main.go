package main

import (
	"github.com/alexflint/go-arg"
	"github.com/pluveto/ankiterm/x/ankicc"
	"github.com/pluveto/ankiterm/x/automata"
	"github.com/pluveto/ankiterm/x/reviewer/oneline"
	"github.com/pluveto/ankiterm/x/reviewer/streamrv"
)

type Args struct {
	BaseURL  string `arg:"-u,--baseURL" help:"Base URL for the server" default:"http://127.0.0.1:8765"`
	Deck     string `arg:"required,positional" help:"Deck name"`
	Reviewer string `arg:"-r,--reviewer" help:"Reviewer name" default:"stream"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	am := automata.NewAutomata(ankicc.Client{BaseURL: args.BaseURL})
	if args.Reviewer == "oneline" {
		oneline.Execute(am, args.Deck)
		return
	}
	streamrv.Execute(am, args.Deck)
}
