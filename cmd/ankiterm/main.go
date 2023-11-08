package main

import (
	"github.com/alexflint/go-arg"
	"github.com/pluveto/ankiterm/x/ankicc"
	"github.com/pluveto/ankiterm/x/automata"
	"github.com/pluveto/ankiterm/x/reviewer/oneline"
	"github.com/pluveto/ankiterm/x/reviewer/streamrv"
)

var Version = "development"

type Args struct {
	BaseURL  string `arg:"-u,--baseURL" help:"Base URL for the server" default:"http://127.0.0.1:8765"`
	Deck     string `arg:"required,positional" help:"Deck name"`
	Reviewer string `arg:"-r,--reviewer" help:"Reviewer name" default:"stream"`
}

func (Args) Version() string {
	return Version
}

func main() {
	var args Args
	arg.MustParse(&args)

	am := automata.NewAutomata(ankicc.Client{BaseURL: args.BaseURL})
	switch args.Reviewer {
	case "oneline":
		oneline.Execute(am, args.Deck)
		return
	case "stream":
		streamrv.Execute(am, args.Deck)
	default:
		panic("unknown reviewer: " + args.Reviewer)
	}
}
